package ice

import (
	"reflect"
	"strings"

	ice "shylinux.com/x/icebergs"
	_ "shylinux.com/x/icebergs/base"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/web"
	_ "shylinux.com/x/icebergs/core"
	_ "shylinux.com/x/icebergs/misc"
	kit "shylinux.com/x/toolkits"
)

type Any = ice.Any
type Map = ice.Map

type Message struct{ *ice.Message }

func (m *Message) Spawn(arg ...Any) *Message {
	return &Message{m.Message.Spawn(arg...)}
}
func (m *Message) PushStream() func() *Message {
	cli.PushStream(m.Message)
	return func() *Message {
		m.StatusTimeCount()
		m.ProcessHold()
		return m
	}
}
func (m *Message) Conf(arg ...Any) string {
	return m.Message.Conf(trans(arg...)...)
}
func (m *Message) Cmd(arg ...Any) *Message {
	return &Message{m.Message.Cmd(trans(arg...)...)}
}
func (m *Message) Cmdx(arg ...Any) string {
	return m.Message.Cmdx(trans(arg...)...)
}
func (m *Message) Cmdy(arg ...Any) *Message {
	return &Message{m.Message.Cmdy(trans(arg...)...)}
}
func (m *Message) HTTP(path string, hand Any) {
	if path == "" {
		path = m.CommandKey()
	}
	if !strings.HasPrefix(path, ice.PS) {
		path = ice.PS + path
	}
	if m.Target().Commands[web.WEB_LOGIN] == nil {
		m.Target().Commands[web.WEB_LOGIN] = &ice.Command{Hand: func(m *ice.Message, arg ...string) {
		}}
	}
	m.Target().Commands[path] = &ice.Command{Name: path, Help: path, Hand: func(m *ice.Message, arg ...string) {
		switch hand := hand.(type) {
		case func(*Message, string, ...string):
			hand(&Message{m}, m.CommandKey(), arg...)
		case func(*Message, ...string):
			hand(&Message{m}, arg...)
		case string:
			m.Cmdy(kit.Select(m.CommandKey(), hand), arg)
		}
	}}
}

var Pulse = ice.Pulse

func Render(m *Message, t string, arg ...Any) string {
	return ice.Render(m.Message, t, arg...)
}
func GetTypeKey(obj Any) string {
	switch t, v := ref(obj); v.Kind() {
	case reflect.Struct:
		return kit.Select(t.String(), listKey(t))
	default:
		return ""
	}
}
func GetTypeName(arg ...Any) string {
	return kit.Slice(kit.Split(GetTypeKey(arg[0]), "."), -1)[0]
}
func trans(arg ...Any) []Any {
	if len(arg) > 1 {
		switch action := arg[1].(type) {
		case []string:
		case string:
		default:
			switch _, v := ref(action); v.Kind() {
			case reflect.Func:
				arg[1] = kit.LowerCapital(kit.FuncName(action))
			}
		}
	}

	switch cmd := arg[0].(type) {
	case []string:
	case string:
	default:
		switch t, v := ref(cmd); v.Kind() {
		case reflect.Struct:
			arg[0] = GetTypeKey(t)
		default:
			return append(kit.List(kit.FileName(cmd), ctx.ACTION, kit.LowerCapital(kit.FuncName(cmd))), arg[1:]...)
		}
	}
	return arg
}
