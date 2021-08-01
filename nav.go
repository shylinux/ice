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
	m.Process("_location", "../")
}
func (n Nav) List(m *Message, arg ...string) {
	if len(arg) > 0 {
		m.Process("_location", arg[0])
	} else {
		m.Option(nfs.DIR_ROOT, path.Join(n.Home, strings.TrimPrefix(path.Dir(m.R.URL.Path), n.Prefix)))
		m.Cmdy(nfs.DIR, arg)
	}
}

func (n Nav) Show(key string, show []*Show) []*Show {
	return append([]*Show{
		{Name: "up", Help: "上一级"},
		{Name: _name(key, -1) + " path auto up", Help: "导航"},
	}, show...)
}
func (n Nav) ShortDef() string { return "" }
func (n Nav) FieldDef() string { return "time,hash,type,name,text" }
