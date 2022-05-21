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

type Message struct{ *ice.Message }

func Render(m *Message, t string, arg ...interface{}) string {
	return ice.Render(m.Message, t, arg...)
}
func (m *Message) Spawn(arg ...interface{}) *Message {
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
func Name(arg ...interface{}) string {
	switch cmd := arg[0].(type) {
	case string:
	default:
		switch t, v := ref(cmd); v.Kind() {
		case reflect.Struct:
			return kit.Slice(kit.Split(kit.Select(t.String(), listKey(t)), "."), -1)[0]
		}
	}
	return ""
}
func trans(arg ...interface{}) []interface{} {
	if len(arg) > 1 {
		switch action := arg[1].(type) {
		case string:
		case []string:
		default:
			switch _, v := ref(action); v.Kind() {
			case reflect.Func:
				arg[1] = strings.ToLower(kit.FuncName(action))
			}
		}
	}

	switch cmd := arg[0].(type) {
	case string:
	case []string:
	default:
		switch t, v := ref(cmd); v.Kind() {
		case reflect.Struct:
			arg[0] = kit.Select(t.String(), listKey(t))
		default:
			return append(kit.List(kit.FileName(cmd), ctx.ACTION, strings.ToLower(kit.FuncName(cmd))), arg[1:]...)
		}
	}
	return arg
}
func (m *Message) Conf(arg ...interface{}) string {
	return m.Message.Conf(trans(arg...)...)
}
func (m *Message) Cmd(arg ...interface{}) *Message {
	return &Message{m.Message.Cmd(trans(arg...)...)}
}
func (m *Message) Cmdx(arg ...interface{}) string {
	return m.Message.Cmdx(trans(arg...)...)
}
func (m *Message) Cmdy(arg ...interface{}) *Message {
	return &Message{m.Message.Cmdy(trans(arg...)...)}
}

func (m *Message) HTTP(path string, hand interface{}) {
	if path == "" {
		path = m.CommandKey()
	}
	if !strings.HasPrefix(path, ice.PS) {
		path = ice.PS + path
	}
	if m.Target().Commands[web.WEB_LOGIN] == nil {
		m.Target().Commands[web.WEB_LOGIN] = &ice.Command{Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
		}}
	}
	m.Target().Commands[path] = &ice.Command{Name: path, Help: path, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
		switch hand := hand.(type) {
		case func(*Message, string, ...string):
			hand(&Message{m}, cmd, arg...)
		case func(*Message, ...string):
			hand(&Message{m}, arg...)
		case string:
			m.Cmdy(kit.Select(m.CommandKey(), hand), arg)
		}
	}}
}

var Pulse = ice.Pulse
