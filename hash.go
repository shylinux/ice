package ice

import (
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

type Hash struct {
	Data

	short string `data:""`
	field string `data:"time,hash,type,name,text"`

	create string `name:"create type name text" help:"创建"`
	remove string `name:"remove" help:"删除"`
	list   string `name:"list hash auto create" help:"缓存"`
}

func (h Hash) Short(m *Message) string {
	return kit.Select(mdb.HASH, h.Data.Short(m))
}

func (h Hash) Prunes(m *Message, arg ...string) {
	h.Data.Prunes(m, mdb.HASH, arg)
}
func (h Hash) Inputs(m *Message, arg ...string) {
	h.Data.Inputs(m, mdb.HASH, arg)
}
func (h Hash) Create(m *Message, arg ...string) *Message {
	h.Data.Insert(m, mdb.HASH, arg)
	return m
}
func (h Hash) Remove(m *Message, arg ...string) {
	h.Data.Delete(m, mdb.HASH, m.OptionSimple(h.Short(m)))
}
func (h Hash) Modify(m *Message, arg ...string) {
	h.Data.Modify(m, mdb.HASH, m.OptionSimple(h.Short(m)), arg)
}
func (h Hash) Export(m *Message, arg ...string) {
	h.Data.Export(m, mdb.HASH)
}
func (h Hash) Import(m *Message, arg ...string) {
	h.Data.Import(m, mdb.HASH)
}
func (h Hash) List(m *Message, arg ...string) *Message {
	m.Fields(len(arg), h.Field(m))
	h.Data.Select(m, mdb.HASH, h.Short(m), arg)
	m.PushAction(h.Remove)
	return m
}
