package ice

import (
	kit "shylinux.com/x/toolkits"
)

func App(index string, text string) *Page {
	index = kit.Keys("web", index)
	p := &Page{name: index, args: map[int]string{}}
	for i, line := range kit.Split(text, "\n", "\n", "\n") {
		list := kit.Split(line)
		p.list = append(p.list, list[0])
		p.args[i] = kit.Format(list[1:])
	}
	Cmd(index, p)
	return p
}

func Arg(arg ...interface{}) string {
	return kit.Format(kit.Simple(arg...))
}
