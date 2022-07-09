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

func ref(obj Any) (reflect.Type, reflect.Value) {
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
func transMethod(obj Any, command *ice.Command, config *ice.Config) {
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

		if key := kit.LowerCapital(t.Method(i).Name); key == mdb.LIST {
			command.Hand = func(m *ice.Message, arg ...string) { h(m, arg...) }
		} else {
			if key == INIT {
				key = CTX_INIT
			}
			if key == EXIT {
				key = CTX_EXIT
			}
			if action, ok := command.Actions[key]; !ok {
				command.Actions[key] = &ice.Action{Hand: h}
			} else {
				action.Hand = h
			}
		}
	}
}
func transField(obj Any, command *ice.Command, config *ice.Config) {
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
		key, tag, val := t.Field(i).Name, t.Field(i).Tag, ""
		if v.Field(i).CanInterface() {
			switch v := v.Field(i).Interface().(type) {
			case string:
				val = v
			}
		}

		if data, ok := tag.Lookup(mdb.DATA); ok { // data tag
			kit.Value(meta, key, kit.Select(data, val))
		}

		if name, ok := tag.Lookup(mdb.NAME); ok { // name tag
			name = kit.Select(name, val)
			if help := tag.Get(mdb.HELP); key == mdb.LIST {
				config.Name, config.Help = name, help
				command.Name, command.Help = name, help
			} else if action, ok := command.Actions[key]; ok {
				action.Name, action.Help = name, help
			}
		}

		if http, ok := tag.Lookup(HTTP); ok { // http tag
			hand := func(msg *Message, arg ...string) { msg.Cmdy(msg.CommandKey(), ctx.ACTION, key, arg) }
			if key == mdb.LIST {
				hand = func(msg *Message, arg ...string) { msg.Cmdy(msg.CommandKey(), arg) }
			}

			last := command.Actions[CTX_INIT]
			command.Actions[CTX_INIT] = &ice.Action{Hand: func(m *ice.Message, arg ...string) {
				if last != nil && last.Hand != nil {
					last.Hand(m, arg...)
				}
				(&Message{m}).HTTP(kit.Select(m.CommandKey(), http), hand)
			}}
		}
	}
}

var list = map[string]string{}

func listKey(t reflect.Type, arg ...string) string {
	if len(arg) == 0 {
		return list[kit.Keys(t.PkgPath(), t.String())]
	}
	list[kit.Keys(t.PkgPath(), t.String())] = arg[0]
	return arg[0]
}
func cmd(key string, obj Any, arg ...Any) string {
	if obj == nil {
		return key
	}

	config := &ice.Config{Value: kit.Data(arg...)}
	command := &ice.Command{Name: mdb.LIST, Help: "列表", Actions: map[string]*ice.Action{}, Meta: kit.Dict()}

	switch obj := obj.(type) {
	case func(*Message, ...string):
		command.Hand = func(m *ice.Message, arg ...string) { obj(&Message{m}, arg...) }

	default:
		t, _ := ref(obj)
		listKey(t, key)
		ice.AddFileCmd(ice.FileRequire(4), key)

		transMethod(obj, command, config)
		transField(obj, command, config)
	}

	if strings.HasPrefix(command.Name, mdb.LIST) {
		command.Name = strings.Replace(command.Name, mdb.LIST, kit.Slice(strings.Split(key, ice.PT), -1)[0], 1)
	}

	last, list := ice.Index, strings.Split(key, ice.PT)
	for i := 1; i < len(list); i++ {
		has := false
		if ice.Pulse.Search(strings.Join(list[:i], ice.PT)+ice.PT, func(p *ice.Context, s *ice.Context) { has, last = true, s }); !has {
			context := &ice.Context{Name: list[i-1]}
			last.Register(context, &web.Frame{})
			last = context
		}

		if i == len(list)-1 {
			last.Merge(&ice.Context{Configs: map[string]*ice.Config{list[i]: config}, Commands: map[string]*ice.Command{list[i]: command}})
		}
	}
	return key
}
func Cmd(key string, obj Any, arg ...Any) string { return cmd(key, obj, arg...) }
