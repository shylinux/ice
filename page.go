package ice

import (
	"strings"

	"shylinux.com/x/icebergs/base/ctx"
	kit "shylinux.com/x/toolkits"
)

type Page struct {
	Tool
	name string
	list []string
	args map[int]string
}

func (p Page) Command(m *Message, arg ...string) {
	if len(arg) == 0 {
		for i, item := range p.list {
			m.Push("index", kit.Select(kit.Keys(p.name, item), item, strings.Contains(item, ".")))
			m.Push("args", kit.Select("[]", p.args[i]))
		}
		return
	}
	m.Cmdy(ctx.COMMAND, arg)
}
func (p Page) Show(show []*Show) []*Show {
	return append([]*Show{
		{Name: "command cmd...", Help: "命令"},
		{Name: "run", Help: "执行"},
		{Name: "list hash auto command", Help: "工具"},
	}, show...)
}
func (p Page) ShortDef() string { return "" }
func (p Page) FieldDef() string { return "time,hash,type,name,text" }

func (p *Page) Cmd(key string, obj interface{}, shows ...[]*Show) *Page {
	switch obj := obj.(type) {
	case string:
		if p.args == nil {
			p.args = map[int]string{}
		}
		p.args[len(p.list)] = obj
	case nil:
	default:
		Cmd(kit.Keys(p.name, key), obj, shows...)
	}
	p.list = append(p.list, key)
	return p
}

func App(ctx, cmd string, cb func(*Page)) *Page {
	p := &Page{name: ctx}
	cb(p)
	Cmd(kit.Keys(ctx, cmd), p)
	return p
}

func Arg(arg ...interface{}) string {
	return kit.Format(kit.Simple(arg...))
}
