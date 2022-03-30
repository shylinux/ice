package ice

import (
	"path"
	"reflect"
	"strings"

	ice "shylinux.com/x/icebergs"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/nfs"
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
	deploy   string `name:"deploy" help:"部署"`
	list     string `name:"list port path auto start order build download" help:"源码"`
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

func (c Code) Download(m *Message, src string, arg ...string) {
	m.Cmdy(code.INSTALL, web.DOWNLOAD, src, arg)
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
func (c Code) Start(m *Message, src string, arg ...string) {
	m.Cmdy(code.INSTALL, cli.START, src, arg)
}
func (c Code) List(m *Message, src string, arg ...string) {
	m.Cmdy(code.INSTALL, src, arg)
}
func (c Code) Trash(m *Message, src string, arg ...string) {
	m.Cmdy(code.INSTALL, nfs.TRASH)
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
