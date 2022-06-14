package ice

import (
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

type Lang struct {
}

func (l Lang) Init(m *Message, arg ...Any) {
	m.Config(code.PLUG, kit.Dict(arg...))
	code.LoadPlug(m.Message, m.PrefixKey())
	m.Config(kit.Keys(code.PLUG, code.PREPARE), "")
	m.Cmd(mdb.PLUGIN, mdb.CREATE, m.PrefixKey())
	m.Cmd(mdb.RENDER, mdb.CREATE, m.PrefixKey())
	m.Cmd(mdb.ENGINE, mdb.CREATE, m.PrefixKey())
	m.Cmd(mdb.SEARCH, mdb.CREATE, m.PrefixKey())
}
func (l Lang) Plugin(m *Message, arg ...string) {
	m.Echo(m.Config(code.PLUG))
}
func (l Lang) Render(m *Message, arg ...string) {
}
func (l Lang) Engine(m *Message, arg ...string) {
}
func (l Lang) Search(m *Message, arg ...string) {
}

func (l Lang) System(m *Message, arg ...Any) bool {
	if !code.InstallSoftware(m.Spawn().Message, kit.Simple(arg)[0], m.Configv(INSTALL)) {
		return false
	}
	if cli.IsSuccess(m.Cmdy(append(kit.List(cli.SYSTEM), arg...)...).Message) {
		m.SetAppend()
		return true
	}
	return false
}
