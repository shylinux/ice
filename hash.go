package ice

import (
	"shylinux.com/x/icebergs/base/mdb"
)

type Hash struct {
	Data

	short string `data:""`
	field string `data:"time,hash,type,name,text"`

	create string `name:"create type=hi name=hello text=worl" help:"创建"`
	list   string `name:"list hash auto create" help:"缓存"`
}

func (h Hash) Inputs(m *Message, arg ...string) { mdb.HashInputs(m.Message, arg) }
func (h Hash) Create(m *Message, arg ...string) { mdb.HashCreate(m.Message, arg) }
func (h Hash) Remove(m *Message, arg ...string) { mdb.HashRemove(m.Message, arg) }
func (h Hash) Modify(m *Message, arg ...string) { mdb.HashModify(m.Message, arg) }
func (h Hash) Export(m *Message, arg ...string) { mdb.HashExport(m.Message, arg) }
func (h Hash) Import(m *Message, arg ...string) { mdb.HashImport(m.Message, arg) }
func (h Hash) List(m *Message, arg ...string) *Message {
	mdb.HashSelect(m.Message, arg...)
	return m
}
