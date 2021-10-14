package ice

import (
	"reflect"
	"strings"

	ice "shylinux.com/x/icebergs"
	"shylinux.com/x/icebergs/base/web"
	kit "shylinux.com/x/toolkits"
)

func ref(obj interface{}) (reflect.Type, reflect.Value) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t, v = t.Elem(), v.Elem()
	}
	return t, v
}
func val(m *ice.Message, arg ...string) []reflect.Value {
	args := []reflect.Value{reflect.ValueOf(&Message{m})}
	for _, v := range arg {
		args = append(args, reflect.ValueOf(v))
	}
	return args
}
func transMethod(config *ice.Config, command *ice.Command, obj interface{}) {
	t, v := ref(obj)
	for i := 0; i < v.NumMethod(); i++ {
		method := v.Method(i)
		var h func(*ice.Message, ...string)
		switch method.Interface().(type) {
		case func(*Message, ...string):
			h = func(m *ice.Message, arg ...string) { method.Call(val(m, arg...)) }
		case func(*Message):
			h = func(m *ice.Message, arg ...string) { method.Call(val(m)) }
		default:
			continue
		}

		if key := strings.ToLower(t.Method(i).Name); key == "list" {
			command.Hand = func(m *ice.Message, c *ice.Context, cmd string, arg ...string) { h(m, arg...) }
		} else {
			if action, ok := command.Action[key]; !ok {
				command.Action[key] = &ice.Action{Hand: h}
			} else {
				action.Hand = h
			}
		}
	}

}
func transField(config *ice.Config, command *ice.Command, obj interface{}) {
	t, v := ref(obj)
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct {
			if v.Field(i).CanInterface() {
				transField(config, command, v.Field(i).Interface())
			}
		}
	}

	meta := kit.Value(config.Value, kit.MDB_META)
	for i := 0; i < v.NumField(); i++ {
		key, tag := t.Field(i).Name, t.Field(i).Tag
		if data := tag.Get("data"); data != "" {
			kit.Value(meta, key, data)
		}

		if name := tag.Get("name"); name != "" {
			if help := tag.Get("help"); key == "list" {
				command.Name, command.Help = name, help
				config.Name, config.Help = name, help
			} else if action, ok := command.Action[key]; ok {
				action.Name, action.Help = name, help
			}
		}
	}
}
func Cmd(key string, obj interface{}) {
	if obj == nil {
		return
	}
	config := &ice.Config{Value: kit.Data()}
	command := &ice.Command{Action: map[string]*ice.Action{}}

	switch obj := obj.(type) {
	case func(*Message, ...string):
		command.Hand = func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			obj(&Message{m}, arg...)
		}
	default:
		transMethod(config, command, obj)
		transField(config, command, obj)
	}

	last := ice.Index
	list := strings.Split(key, ".")
	for i := 1; i < len(list); i++ {
		has := false
		ice.Pulse.Search(strings.Join(list[:i], ".")+".", func(p *ice.Context, s *ice.Context) {
			has, last = true, s
		})
		if !has {
			context := &ice.Context{Name: list[i-1]}
			last.Register(context, &web.Frame{})
			last = context
		}
		if i < len(list)-1 {
			continue
		}

		last.Merge(&ice.Context{Commands: map[string]*ice.Command{
			ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
				if action, ok := command.Action[ice.INIT]; ok {
					action.Hand(m, arg...)
				} else {
					m.Load()
				}
			}},
			ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
				m.Save()
			}},
			list[i]: command,
		}, Configs: map[string]*ice.Config{list[i]: config}})
	}
}
