package ice

import (
	ice "shylinux.com/x/icebergs"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/web"
	"shylinux.com/x/icebergs/core/chat"
	kit "shylinux.com/x/toolkits"
)

var prepare []interface{}

func init() {
	Cmd("web.chat.prepare._init", func(m *Message, arg ...string) {
		for _, p := range prepare {
			switch p := p.(type) {
			case func(*Message, ...string):
				p(m, arg...)
			}
		}
	})
}
func App(path, name, text string, arg ...string) {
	prepare = append(prepare, func(m *Message, arg ...string) {
		m.Cmdy(chat.WEBSITE, mdb.CREATE, kit.SimpleKV("path,type,name,text", path, "txt", name, text), arg)
	})
}

func Run(arg ...string) string {
	ice.Pulse.Set(ice.MSG_DETAIL)
	ice.Pulse.Set(ice.MSG_APPEND)
	ice.Pulse.Set(ice.MSG_RESULT)
	return ice.Run(arg...)
}
func RunServe(arg ...string) string {
	return ice.Run(kit.Simple(web.SERVE, cli.START, ice.DEV, "", arg)...)
}
