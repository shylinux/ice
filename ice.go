package ice

import (
	"strings"

	ice "shylinux.com/x/icebergs"
	_ "shylinux.com/x/icebergs/base"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/base/web"
	_ "shylinux.com/x/icebergs/core"
	_ "shylinux.com/x/icebergs/misc"
	kit "shylinux.com/x/toolkits"
)

type Message struct{ *ice.Message }

func (m *Message) Spawn() *Message {
	return &Message{m.Message.Spawn()}
}
func (m *Message) HTTP(path string, hand interface{}) {
	if path == "" {
		path = m.CommandKey()
	}
	if !strings.HasPrefix(path, ice.PS) {
		path = ice.PS + path
	}
	m.Target().Commands[path] = &ice.Command{Name: path, Help: "", Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
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
func (m *Message) Cmdy(arg ...interface{}) *Message {
	switch cmd := arg[0].(type) {
	case string:
	default:
		m.Debug("what %v", arg)
		return &Message{m.Message.Cmdy(kit.FileName(cmd), ctx.ACTION, strings.ToLower(kit.FuncName(cmd)), arg[1:])}
	}
	return &Message{m.Message.Cmdy(arg...)}
}

func Run(arg ...string) string {
	ice.Pulse.Set(ice.MSG_DETAIL)
	ice.Pulse.Set(ice.MSG_APPEND)
	ice.Pulse.Set(ice.MSG_RESULT)
	return ice.Run(arg...)
}
func RunServe(port string, arg ...string) string {
	return ice.Run(kit.Simple(web.SERVE, cli.START, ice.DEV, "", tcp.PORT, port, arg)...)
}
