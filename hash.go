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
	list   string `name:"list hash auto create" help:"缓存"`
}

func (h Hash) Short(m *Message) string {
	return kit.Select(kit.MDB_HASH, h.Data.Short(m))
}

func (h Hash) Inputs(m *Message, arg ...string) {
	h.Data.Inputs(m, mdb.HASH, arg)
}
func (h Hash) Create(m *Message, arg ...string) {
	h.Data.Insert(m, mdb.HASH, arg)
}
func (h Hash) Remove(m *Message, arg ...string) {
	h.Data.Delete(m, mdb.HASH, h.Short(m), m.Option(h.Short(m)))
}
func (h Hash) Modify(m *Message, arg ...string) {
	h.Data.Modify(m, mdb.HASH, h.Short(m), m.Option(h.Short(m)), arg)
}
func (h Hash) List(m *Message, arg ...string) {
	m.Fields(len(arg), h.Field(m))
	h.Data.Select(m, mdb.HASH, h.Short(m), arg)
}
