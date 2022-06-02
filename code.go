package ice

import (
	"path"
	"reflect"
	"runtime"
	"strings"

	ice "shylinux.com/x/icebergs"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/base/web"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

type Code struct {
	serve    string `name:"serve" help:"服务"`
	inputs   string `name:"inputs" help:"补全"`
	download string `name:"download" help:"下载"`
	install  string `name:"install" help:"安装"`
	build    string `name:"build" help:"构建"`
	order    string `name:"order" help:"定制"`
	start    string `name:"start" help:"启动"`
	stop     string `name:"stop" help:"停止"`
	open     string `name:"open" help:"打开"`
	deploy   string `name:"deploy" help:"部署"`
	list     string `name:"list port path auto start order build download" help:"源码"`
}

func (c Code) Link(m *Message, arg ...string) string {
	return kit.Select(m.Config(nfs.SOURCE), kit.Select(m.Config(runtime.GOOS), arg, 0))
}
func (c Code) Path(m *Message, url string) string {
	return path.Join(m.Conf(code.INSTALL, kit.Keym(nfs.PATH)), kit.TrimExt(url))
}
func (c Code) PathOther(m *Message, url string) string {
	p := path.Join(m.Conf(code.INSTALL, kit.Keym(nfs.PATH)), path.Base(url))
	return kit.Path(m.Conf(code.INSTALL, kit.Keym(nfs.PATH)), strings.Split(m.Cmdx(cli.SYSTEM, "sh", "-c", kit.Format("tar tf %s| head -n1", p)), "/")[0])
}
func (c Code) Prepare(m *Message, cb interface{}) {
	m.Optionv(code.PREPARE, cb)
}

func (c Code) Inputs(m *Message, arg ...string) {
	switch arg[0] {
	case tcp.PORT:
		m.Cmdy(tcp.PORT)
	case tcp.HOST:
		m.Cmdy(tcp.HOST).Cut(aaa.IP).RenameAppend(aaa.IP, tcp.HOST)
	}
}
func (c Code) Install(m *Message, arg ...string) {
	m.Cmdy(code.INSTALL, web.DOWNLOAD, c.Link(m), kit.Slice(arg, 1))
}
func (c Code) Download(m *Message, arg ...string) {
	m.Cmdy(code.INSTALL, web.DOWNLOAD, kit.Select(m.Config(nfs.SOURCE), arg, 0), kit.Slice(arg, 1))
}
func (c Code) Source(m *Message, src string, arg ...string) {
	m.Cmdy(code.INSTALL, nfs.SOURCE, src, arg)
}
func (c Code) Build(m *Message, src string, arg ...string) {
	m.Cmdy(code.INSTALL, cli.BUILD, src, arg)
}
func (c Code) Order(m *Message, src, dir string, arg ...string) {
	m.Cmd(nfs.PUSH, ice.ETC_PATH, kit.Path(m.Conf(code.INSTALL, kit.Keym(nfs.PATH)), kit.TrimExt(src), dir+ice.NL))
	m.Cmdy(nfs.CAT, ice.ETC_PATH)
}
func (c Code) Start(m *Message, src string, arg ...interface{}) {
	args := []string{}
	for _, v := range arg {
		switch v := v.(type) {
		case func(string) []string:
			m.Option(code.PREPARE, v)
		case string:
			args = append(args, v)
		default:
			args = append(args, kit.Format(v))
		}
	}
	m.Cmdy(code.INSTALL, cli.START, c.Link(m, src), args)
}
func (c Code) Stop(m *Message, arg ...string) {
	m.Cmdy(code.INSTALL, cli.STOP, arg)
}
func (c Code) Open(m *Message, arg ...string) {
	m.ProcessOpen(m.Option(mdb.LINK))
}
func (c Code) List(m *Message, src string, arg ...string) {
	m.Cmdy(code.INSTALL, c.Link(m, src), arg)
}
func (c Code) Trash(m *Message, arg ...string) {
	m.Cmdy(nfs.TRASH, m.Option(nfs.PATH))
}
func (c Code) PushLink(m *Message) *Message {
	hostname := kit.ParseURLMap(m.Option(ice.MSG_USERWEB))[tcp.HOSTNAME]
	m.Tables(func(value map[string]string) { m.Push(mdb.LINK, kit.Format("http://%s:%s", hostname, value[tcp.PORT])) })
	return m
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
	m.PushStream()
	c.System(m, dir, arg...)
	m.ProcessHold()
	m.StatusTime()
}

func CodeCmd(obj interface{}, arg ...interface{}) string {
	return cmd(kit.Keys("web.code", kit.FileName(2)), obj, arg...)
}

func modName(str string) string {
	ls := strings.Split(str, ice.PS)
	mod := ls[0]
	if strings.Contains(ls[0], ice.PT) {
		mod = kit.Select(mod, ls, 2)
	}
	if strings.HasPrefix(mod, "20") {
		mod = kit.Select(mod, strings.Split(mod, "-"), 1)
	}
	return strings.TrimSuffix(mod, "-story")
}
func getModCmd(p string, n int, obj interface{}) string {
	switch t, v := ref(obj); v.Kind() {
	case reflect.Struct:
		return kit.Keys(p, modName(t.PkgPath()), strings.ToLower(kit.Slice(kit.Split(t.String(), ice.PT), -1)[0]))
	default:
		return kit.Keys(p, modName(kit.ModName(n+1)), kit.FileName(n+1))
	}
}
func CodeModCmd(obj interface{}, arg ...interface{}) string {
	return cmd(getModCmd("web.code", 2, obj), obj, arg...)
}

func ctxName(str string) string {
	ls := strings.Split(str, ice.PS)
	if strings.Contains(ls[0], ice.PT) {
		ls = kit.Slice(ls, 3)
	} else {
		ls = kit.Slice(ls, 1)
	}
	if ls[0] == ice.SRC {
		ls = kit.Slice(ls, 1)
	}
	if ls[0] == ice.MISC {
		ls = kit.Slice(ls, 1)
	}
	return strings.Join(ls, ice.PT)
}
func getCtxCmd(p string, n int, obj interface{}) string {
	switch t, v := ref(obj); v.Kind() {
	case reflect.Struct:
		return kit.Keys(p, ctxName(t.PkgPath()), strings.ToLower(kit.Slice(kit.Split(t.String(), ice.PT), -1)[0]))
	default:
		return kit.Keys(p, kit.PathName(n+1), kit.FileName(n+1))
	}
}
func CodeCtxCmd(obj interface{}, arg ...interface{}) string {
	return cmd(getCtxCmd("web.code", 2, obj), obj, arg...)
}
