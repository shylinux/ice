package ice

import (
	"path"
	"strings"

	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/web"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

type Code struct {
	inputs   string `name:"inputs" help:"补全"`
	download string `name:"download" help:"下载"`
	build    string `name:"build" help:"构建"`
	order    string `name:"order" help:"加载"`
	start    string `name:"start" help:"启动"`
	list     string `name:"list port path auto start build download" help:"服务器"`
}

func (c Code) Path(m *Message, url string) string {
	return path.Join(m.Conf(code.INSTALL, kit.META_PATH), kit.TrimExt(url))
}
func (c Code) PathOther(m *Message, url string) string {
	p := path.Join(m.Conf(code.INSTALL, kit.META_PATH), path.Base(url))
	return kit.Path(m.Conf(code.INSTALL, kit.META_PATH), strings.Split(m.Cmdx(cli.SYSTEM, "sh", "-c", kit.Format("tar tf %s| head -n1", p)), "/")[0])
}
func (c Code) Prepare(m *Message, cb interface{}) {
	m.Optionv(code.PREPARE, cb)
}

func (c Code) Download(m *Message, src string, arg ...string) {
	m.Cmdy(code.INSTALL, web.DOWNLOAD, src, arg)
}
func (c Code) Build(m *Message, src string, arg ...string) {
	m.Cmdy(code.INSTALL, cli.BUILD, src, arg)
}
func (c Code) Order(m *Message, src, dir string, arg ...string) {
	m.Cmd(nfs.PUSH, "etc/path", kit.Path(m.Conf(code.INSTALL, kit.META_PATH), kit.TrimExt(src), dir+"\n"))
	m.Cmdy(nfs.CAT, "etc/path")
}
func (c Code) Start(m *Message, src string, arg ...string) {
	m.Cmdy(code.INSTALL, cli.START, src, arg)
}
func (c Code) List(m *Message, src string, arg ...string) {
	m.Cmdy(code.INSTALL, src, arg)
}

func (c Code) Daemon(m *Message, dir string, arg ...string) {
	m.Option(cli.CMD_DIR, dir)
	m.Cmdy(cli.DAEMON, arg)
}
func (c Code) System(m *Message, dir string, arg ...string) {
	m.Option(cli.CMD_DIR, dir)
	m.Cmdy(cli.SYSTEM, arg)
}
func (c Code) Stream(m *Message, dir string, arg ...string) {
	web.PushStream(m.Message)
	c.System(m, dir, arg...)
	m.ProcessHold()
	m.StatusTime()
}

func WikiCmd(obj interface{})    { Cmd(kit.Keys("web.wiki", kit.FileName(2)), obj) }
func CodeCmd(obj interface{})    { Cmd(kit.Keys("web.code", kit.FileName(2)), obj) }
func CodeCtxCmd(obj interface{}) { Cmd(kit.Keys("web.code", kit.PathName(2), kit.FileName(2)), obj) }
func CodeModCmd(obj interface{}) {
	Cmd(kit.Keys("web.code", strings.TrimSuffix(kit.ModName(obj), "-story"), kit.FileName(2)), obj)
}
