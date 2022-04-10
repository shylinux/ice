package ice

import (
	"reflect"
	"strings"

	ice "shylinux.com/x/icebergs"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/mdb"
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
func transMethod(obj interface{}, command *ice.Command, config *ice.Config) {
	t, v := ref(obj)
	for i := 0; i < v.NumMethod(); i++ {
		method := v.Method(i)

		var h func(*ice.Message, ...string)
		switch method.Interface().(type) {
		case func(*Message, ...string) *Message:
			h = func(m *ice.Message, arg ...string) { method.Call(val(m, arg...)) }
		case func(*Message, ...string):
			h = func(m *ice.Message, arg ...string) { method.Call(val(m, arg...)) }
		case func(*Message):
			h = func(m *ice.Message, arg ...string) { method.Call(val(m)) }
		default:
			continue
		}

		if key := strings.ToLower(t.Method(i).Name); key == mdb.LIST {
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

func transField(obj interface{}, command *ice.Command, config *ice.Config) {
	t, v := ref(obj)
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct {
			if v.Field(i).CanInterface() {
				transField(v.Field(i).Interface(), command, config)
			}
		}
	}

	meta := kit.Value(config.Value, mdb.META)
	for i := 0; i < v.NumField(); i++ {
		key, tag := t.Field(i).Name, t.Field(i).Tag
		switch key {
		case "display":
			for k, v := range ice.DisplayRequire(3, tag.Get(mdb.DATA)) {
				command.Meta[k] = v
			}
			continue
		}
		if data := tag.Get(mdb.DATA); data != "" {
			kit.Value(meta, key, data)
		}

		name := tag.Get(mdb.NAME)
		if name == "" {
			continue
		}

		if help := tag.Get(mdb.HELP); key == mdb.LIST {
			config.Name, config.Help = name, help
			command.Name, command.Help = name, help
		} else if action, ok := command.Action[key]; ok {
			action.Name, action.Help = name, help
		}

		h, ok := tag.Lookup(ice.HTTP)
		if !ok {
			continue
		}

		hand := func(msg *Message, arg ...string) { msg.Cmdy(msg.CommandKey(), ctx.ACTION, key, arg) }
		if key == mdb.LIST {
			hand = func(msg *Message, arg ...string) { msg.Cmdy(msg.CommandKey(), arg) }
		}
		last := command.Action[ice.CTX_INIT]
		command.Action[ice.CTX_INIT] = &ice.Action{Hand: func(m *ice.Message, arg ...string) {
			if last != nil && last.Hand != nil {
				last.Hand(m, arg...)
			}
			(&Message{m}).HTTP(kit.Select(m.CommandKey(), h), hand)
		}}
	}
}

var list = map[string]string{}

func GetTypeKey(obj interface{}) string {
	switch t, v := ref(obj); v.Kind() {
	case reflect.Struct:
		return kit.Select(t.String(), listKey(t))
	default:
		return ""
	}
}
func listKey(t reflect.Type, arg ...string) string {
	if len(arg) == 0 {
		return list[kit.Keys(t.PkgPath(), t.String())]
	}
	list[kit.Keys(t.PkgPath(), t.String())] = arg[0]
	return arg[0]
}

func Cmd(key string, obj interface{}, arg ...interface{}) string { return cmd(key, obj, arg...) }
func cmd(key string, obj interface{}, arg ...interface{}) string {
	if obj == nil {
		return key
	}
	command := &ice.Command{Name: mdb.LIST, Help: "列表", Action: map[string]*ice.Action{}, Meta: kit.Dict()}
	config := &ice.Config{Value: kit.Data(arg...)}

	switch obj := obj.(type) {
	case func(*Message, ...string):
		command.Hand = func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			obj(&Message{m}, arg...)
		}
	default:
		t, _ := ref(obj)
		listKey(t, key)
		p := kit.FileLine(3, 100)
		ice.AddFileCmd(p, key)

		transMethod(obj, command, config)
		transField(obj, command, config)
	}
	if strings.HasPrefix(command.Name, mdb.LIST) {
		command.Name = strings.Replace(command.Name, mdb.LIST, kit.Slice(strings.Split(key, ice.PT), -1)[0], 1)
	}

	last := ice.Index
	list := strings.Split(key, ice.PT)
	for i := 1; i < len(list); i++ {
		has := false
		ice.Pulse.Search(strings.Join(list[:i], ice.PT)+ice.PT, func(p *ice.Context, s *ice.Context) {
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
					action.Hand(m.Spawn(kit.Ext(key)), arg...)
				}
			}},
			ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
				if action, ok := command.Action[ice.EXIT]; ok {
					action.Hand(m.Spawn(kit.Ext(key)), arg...)
				}
			}},
			list[i]: command,
		}, Configs: map[string]*ice.Config{list[i]: config}})
	}
	return key
}
