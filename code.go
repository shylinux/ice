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
	inputs     string `name:"inputs" help:"补全"`
	install    string `name:"install" help:"安装"`
	download   string `name:"download" help:"下载"`
	build      string `name:"build" help:"构建"`
	order      string `name:"order" help:"定制"`
	start      string `name:"start port" help:"启动"`
	stop       string `name:"stop" help:"停止"`
	open       string `name:"open" help:"打开"`
	serve      string `name:"serve" help:"服务"`
	deploy     string `name:"deploy" help:"部署"`
	list       string `name:"list port path auto start build download" help:"源代码"`
	listScript string `name:"listScript" help:"脚本"`
	runScript  string `name:"runScript" help:"执行"`
	catScript  string `name:"catScript" help:"查看"`
}

func (s Code) Link(m *Message, arg ...string) string {
	return kit.Select(m.Config(nfs.SOURCE), kit.Select(m.Config(runtime.GOOS), arg, 0))
}
func (s Code) Path(m *Message, url string) string {
	return path.Join(ice.USR_INSTALL, kit.TrimExt(url))
}
func (s Code) PathOther(m *Message, url string) string {
	return kit.Path(ice.USR_INSTALL, strings.Split(m.Cmdx(cli.SYSTEM, "sh", "-c", kit.Format("tar tf %s| head -n1", path.Join(ice.USR_INSTALL, path.Base(url)))), "/")[0])
}

func (s Code) Inputs(m *Message, arg ...string) {
	switch arg[0] {
	case tcp.PORT:
		m.Cmdy(tcp.PORT)
	case tcp.HOST:
		m.Cmdy(tcp.HOST).Cut(aaa.IP).RenameAppend(aaa.IP, tcp.HOST)
	}
}
func (s Code) Install(m *Message, arg ...string) {
	m.Cmdy(code.INSTALL, web.DOWNLOAD, s.Link(m, kit.Select("", arg, 0)), kit.Slice(arg, 1))
}
func (s Code) Download(m *Message, arg ...string) {
	m.Cmdy(code.INSTALL, web.DOWNLOAD, kit.Select(m.Config(nfs.SOURCE), arg, 0), kit.Slice(arg, 1))
}
func (s Code) Source(m *Message, src string, arg ...string) {
	m.Cmdy(code.INSTALL, nfs.SOURCE, s.Link(m, src), kit.Select(nfs.PWD, arg, 0))
}
func (s Code) prepare(m *Message, arg ...Any) []string {
	args := []string{}
	for _, v := range arg {
		switch v := v.(type) {
		case func(string) []string:
			m.Option(code.PREPARE, v)
		case func(string):
			m.Option(code.PREPARE, v)
		case []string:
			args = append(args, v...)
		case string:
			args = append(args, v)
		default:
			args = append(args, kit.Format(v))
		}
	}
	return args
}
func (s Code) Build(m *Message, arg ...Any) {
	m.PushStream()
	args := s.prepare(m, arg...)
	m.Cmdy(code.INSTALL, cli.BUILD, kit.Select(m.Config(nfs.SOURCE), args, 0), kit.Slice(args, 1))
}
func (s Code) Order(m *Message, src, dir string, arg ...string) {
	m.Cmdy(code.INSTALL, cli.ORDER, s.Link(m, src), dir)
}
func (s Code) Start(m *Message, src, bin string, arg ...Any) {
	args := s.prepare(m, arg...)
	m.Cmdy(code.INSTALL, cli.START, s.Link(m, src), bin, args)
}
func (s Code) Stop(m *Message, arg ...string) {
	m.Cmdy(code.INSTALL, cli.STOP, arg)
}
func (s Code) Open(m *Message, arg ...string) {
	m.ProcessOpen(m.Option(mdb.LINK))
}
func (s Code) List(m *Message, src string, arg ...string) {
	if m.Cmdy(code.INSTALL, s.Link(m, src), arg); len(arg) == 0 {
		s.PushLink(m)
	}
}
func (s Code) PushLink(m *Message) *Message {
	hostname := m.OptionUserWeb().Hostname()
	m.Tables(func(value map[string]string) { m.Push(mdb.LINK, kit.Format("http://%s:%s", hostname, value[tcp.PORT])) })
	return m
}
func (s Code) ListScript(m *Message) {
	m.Option(nfs.DIR_REG, m.Config("regexp"))
	m.Option(nfs.DIR_DEEP, ice.TRUE)
	m.Cmdy(nfs.DIR, ice.SRC)
	m.PushAction(s.CatScript, s.RunScript)
}
func (s Code) CatScript(m *Message) {
	m.Cmdy(nfs.CAT, m.Option(nfs.PATH))
}
func (s Code) RunScript(m *Message) {
	s.System(m, nfs.PWD, kit.Simple(kit.Split(m.Config("command")), m.Option(nfs.PATH))...)
}

func (s Code) Daemon(m *Message, dir string, arg ...string) {
	m.Option(cli.CMD_DIR, dir)
	if m.Cmdy(cli.DAEMON, arg); cli.IsSuccess(m.Message) {
		m.SetAppend()
	}
}
func (s Code) System(m *Message, dir string, arg ...string) {
	m.Option(cli.CMD_DIR, dir)
	if m.Cmdy(cli.SYSTEM, arg); cli.IsSuccess(m.Message) {
		m.SetAppend()
	}
}
func (s Code) Stream(m *Message, dir string, arg ...string) {
	m.PushStream()
	s.System(m, dir, arg...)
	m.ProcessHold()
	m.StatusTime()
}

func CodeCmd(obj Any, arg ...Any) string {
	return cmd(kit.Keys("web.code", kit.FileName(2)), obj, arg...)
}
func CodeModCmd(obj Any, arg ...Any) string {
	return cmd(getModCmd("web.code", 2, obj), obj, arg...)
}
func CodeCtxCmd(obj Any, arg ...Any) string {
	return cmd(getCtxCmd("web.code", 2, obj), obj, arg...)
}

func getModCmd(p string, n int, obj Any) string {
	switch t, v := ref(obj); v.Kind() {
	case reflect.Struct:
		return kit.Keys(p, modName(t.PkgPath()), strings.ToLower(kit.Slice(kit.Split(t.String(), ice.PT), -1)[0]))
	default:
		return kit.Keys(p, modName(kit.ModName(n+1)), kit.FileName(n+1))
	}
}
func getCtxCmd(p string, n int, obj Any) string {
	switch t, v := ref(obj); v.Kind() {
	case reflect.Struct:
		return kit.Keys(p, ctxName(t.PkgPath()), strings.ToLower(kit.Slice(kit.Split(t.String(), ice.PT), -1)[0]))
	default:
		return kit.Keys(p, kit.PathName(n+1), kit.FileName(n+1))
	}
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
