package ice

import (
	"shylinux.com/x/icebergs/base/mdb"
)

type Lists struct {
	Data
	field string `data:"time,id,type,name,text"`
}

func (l Lists) Insert(m *Message, arg ...string) *Message {
	l.Data.Insert(m, mdb.LIST, arg)
	return m
}
func (l Lists) Delete(m *Message, arg ...string) {
	l.Data.Delete(m, mdb.LIST, m.Option(mdb.ID))
}
func (l Lists) List(m *Message, arg ...string) {
	m.Fields(len(arg), l.Field(m))
	l.Data.Select(m, mdb.LIST, mdb.ID, arg)
}
