package ice

import (
	"path"
	"strings"

	"shylinux.com/x/icebergs/base/nfs"
)

type Nav struct {
	Tool

	Home   string
	Prefix string

	up   string `name:"up" help:"上一级"`
	list string `name:"path auto up" help:"导航"`
}

func (n Nav) Up(m *Message, arg ...string) {
	if strings.TrimPrefix(m.R.URL.Path, n.Prefix) == "/" {
		n.List(m)
		return
	}
	if strings.HasSuffix(m.R.URL.Path, "/") {
		m.ProcessLocation("../")
	} else {
		m.ProcessLocation("./")
	}
}
func (n Nav) List(m *Message, arg ...string) {
	if len(arg) > 0 {
		m.ProcessLocation(arg[0])
		return
	}
	m.Option(nfs.DIR_ROOT, path.Join(n.Home, strings.TrimPrefix(path.Dir(m.R.URL.Path), n.Prefix)))
	m.Cmdy(nfs.DIR, arg)
}
