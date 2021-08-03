package ice

import (
	"path"
	"strings"

	"github.com/shylinux/icebergs/base/nfs"
)

type Nav struct {
	Data
	Home   string
	Prefix string
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

func (n Nav) Show(show []*Show) []*Show {
	return append([]*Show{
		{Name: "up", Help: "上一级"},
		{Name: "list path auto up", Help: "导航"},
	}, show...)
}
func (n Nav) ShortDef() string { return "" }
func (n Nav) FieldDef() string { return "time,hash,type,name,text" }
