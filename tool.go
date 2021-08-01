package ice

import (
	"path"
	"strings"

	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/ctx"
	"github.com/shylinux/icebergs/base/web"
	kit "github.com/shylinux/toolkits"
)

type Tool struct {
	Data
	HTML string
	Home string
}

func (t Tool) Push(m *Message, cmd string, arg ...string) {
	m.Push("index", cmd)
	m.Push("args", kit.Format(arg))
}
func (t Tool) Command(m *Message, arg ...string) {
	if len(arg) == 0 {
		m.Cmd(ctx.COMMAND).Table(func(index int, value map[string]string, head []string) {
			m.Push("index", m.Prefix(value["key"]))
			m.Push("args", kit.Format(kit.Simple()))
		})
		return
	}
	m.Cmdy(ctx.COMMAND, arg)
}
func (t Tool) Run(m *Message, arg ...string) {
	m.Cmdy(arg)
}
func (t Tool) List(m *Message, arg ...string) {
	if strings.HasSuffix(m.R.URL.Path, "/") {
		m.RenderDownload(kit.Select(path.Join(m.Conf(web.SERVE, kit.Keym(ice.VOLCANOS, kit.MDB_PATH)), "page/cmd.html"), path.Join(t.Home, "cmd.html")))
		return
	}
	m.RenderDownload(path.Join(kit.Select(m.Conf(web.SERVE, kit.Keym(ice.VOLCANOS, kit.MDB_PATH)), t.Home), path.Join(arg...)))
}
func (t Tool) Show(key string, show []*Show) []*Show {
	return append([]*Show{
		{Name: "command cmd...", Help: "命令"},
		{Name: "run", Help: "执行"},
		{Name: _name(key, -1) + " hash auto command", Help: "工具"},
	}, show...)
}
func (t Tool) ShortDef() string { return "" }
func (t Tool) FieldDef() string { return "time,hash,type,name,text" }
