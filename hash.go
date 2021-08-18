package ice

import (
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

type Hash struct{ Data }

func (h Hash) Short(m *Message) string {
	return kit.Select(kit.MDB_HASH, h.Data.Short(m))
}
func (h Hash) Field(m *Message) string {
	return kit.Select(h.FieldDef(), h.Data.Field(m))
}

func (h Hash) Create(m *Message, arg ...string) {
	h.Data.Insert(m, mdb.HASH, arg)
}
func (h Hash) Modify(m *Message, arg ...string) {
	h.Data.Modify(m, mdb.HASH, h.Short(m), m.Option(h.Short(m)), arg)
}
func (h Hash) Remove(m *Message, arg ...string) {
	h.Data.Delete(m, mdb.HASH, h.Short(m), m.Option(h.Short(m)))
}
func (h Hash) Inputs(m *Message, arg ...string) {
	h.Data.Inputs(m, mdb.HASH, arg)
}
func (h Hash) List(m *Message, arg ...string) {
	m.Fields(len(arg), h.Field(m))
	h.Data.Select(m, mdb.HASH, h.Short(m), arg)
}

func (h Hash) Show(show []*Show) []*Show {
	return append([]*Show{
		{Name: "create type name text", Help: "创建"},
		{Name: "modify", Help: "编辑"},
		{Name: "remove", Help: "删除"},
		{Name: "inputs", Help: "补全"},
		{Name: "list hash auto create", Help: "缓存"},
	}, show...)
}
func (h Hash) ShortDef() string { return "" }
func (h Hash) FieldDef() string { return "time,hash,type,name,text" }
