package ice

import (
	"path"
	"strings"

	ice "shylinux.com/x/icebergs"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/web"
	kit "shylinux.com/x/toolkits"
)

type Tool struct {
	Home string
}

func (t Tool) Command(m *Message, arg ...string) {
	if len(arg) == 0 {
		m.Cmd(ctx.COMMAND).Table(func(index int, value map[string]string, head []string) {
			m.Push(kit.MDB_INDEX, m.Prefix(value[kit.MDB_KEY]))
			m.Push(kit.MDB_ARGS, kit.Format(kit.Simple()))
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
		m.RenderDownload(path.Join(m.Conf(web.SERVE, kit.Keym(ice.VOLCANOS, kit.MDB_PATH)), "page/cmd.html"))
		return
	}
	m.RenderDownload(path.Join(t.Home, path.Join(arg...)))
}
