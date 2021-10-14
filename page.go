package ice

import (
	"shylinux.com/x/icebergs/base/ctx"
	kit "shylinux.com/x/toolkits"
)

type Page struct {
	Tool
	name string
	list []string
	args map[int]string
}

func (p *Page) Cmd(key string, obj interface{}) *Page {
	switch obj := obj.(type) {
	case string:
		if p.args == nil {
			p.args = map[int]string{}
		}
		p.args[len(p.list)] = obj
	default:
		key = kit.Keys(p.name, key)
		Cmd(key, obj)
	}
	p.list = append(p.list, key)
	return p
}

func (p Page) Command(m *Message, arg ...string) {
	if len(arg) == 0 {
		for i, _ := range p.list {
			m.Push(kit.MDB_INDEX, i)
			m.Push(kit.MDB_ARGS, kit.Select("[]", p.args[i]))
		}
		return
	}
	for _, v := range arg {
		if i := kit.Int(v); i >= 0 && i < len(p.list) {
			m.Cmdy(ctx.COMMAND, p.list[i])
		}
	}
}
func (p Page) Run(m *Message, arg ...string) {
	if i := kit.Int(arg[0]); i >= 0 && i < len(p.list) {
		m.Cmdy(p.list[i], arg[1:])
	}
}
