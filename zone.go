package ice

import (
	"github.com/shylinux/icebergs/base/mdb"
	kit "github.com/shylinux/toolkits"
)

type Zone struct{ Data }

func (z Zone) Short(m *Message) string {
	return kit.Select(z.ShortDef(), z.Data.Short(m))
}
func (z Zone) Field(m *Message) string {
	return kit.Select(z.FieldDef(), z.Data.Field(m))
}

func (z Zone) Create(m *Message, arg ...string) {
	z.Data.Insert(m, mdb.HASH, arg)
}
func (z Zone) Insert(m *Message, arg ...string) {
	z.Data.Insert(m, mdb.HASH, z.Short(m), m.Option(z.Short(m)))
	z.Data.Insert(m, mdb.ZONE, m.Option(z.Short(m)), arg[2:])
}
func (z Zone) Modify(m *Message, arg ...string) {
	z.Data.Modify(m, mdb.ZONE, m.Option(z.Short(m)), m.Option(kit.MDB_ID), arg)
}
func (z Zone) Remove(m *Message, arg ...string) {
	z.Data.Delete(m, mdb.HASH, z.Short(m), m.Option(z.Short(m)))
}
func (z Zone) Inputs(m *Message, arg ...string) {
	if kit.Select("", arg, 0) == z.Short(m) {
		z.Data.Inputs(m, mdb.HASH, arg)
	} else {
		z.Data.Inputs(m, mdb.ZONE, m.Option(z.Short(m)), arg)
	}
}
func (z Zone) List(m *Message, arg ...string) {
	m.Fields(len(arg), "time,"+z.Short(m)+",count", z.Field(m))
	z.Data.Select(m, mdb.ZONE, arg)
}
func (z Zone) Show(show []*Show) []*Show {
	return append([]*Show{
		{Name: "create zone", Help: "创建"},
		{Name: "insert zone type name text", Help: "添加"},
		{Name: "modify", Help: "编辑"},
		{Name: "remove", Help: "删除"},
		{Name: "inputs", Help: "补全"},
		{Name: "list zone id auto insert", Help: "存储"},
	}, show...)
}
func (z Zone) ShortDef() string { return kit.MDB_ZONE }
func (z Zone) FieldDef() string { return "time,id,type,name,text" }
