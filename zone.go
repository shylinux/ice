package ice

import (
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

type Zone struct {
	Data

	short string `data:"zone"`
	field string `data:"time,id,type,name,text"`

	create string `name:"create zone" help:"创建"`
	insert string `name:"insert zone type name text" help:"添加"`
	list   string `name:"list zone id auto insert" help:"存储"`
}

func (z Zone) Inputs(m *Message, arg ...string) {
	if kit.Select("", arg, 0) == z.Short(m) {
		z.Data.Inputs(m, mdb.HASH, arg)
	} else {
		z.Data.Inputs(m, mdb.ZONE, m.Option(z.Short(m)), arg)
	}
}
func (z Zone) Create(m *Message, arg ...string) {
	z.Data.Insert(m, mdb.HASH, arg)
}
func (z Zone) Remove(m *Message, arg ...string) {
	z.Data.Delete(m, mdb.HASH, m.OptionSimple(z.Short(m)))
}
func (z Zone) Insert(m *Message, arg ...string) {
	z.Data.Insert(m, mdb.HASH, m.OptionSimple(z.Short(m)))
	z.Data.Insert(m, mdb.ZONE, m.Option(z.Short(m)), arg[2:])
}
func (z Zone) Modify(m *Message, arg ...string) {
	z.Data.Modify(m, mdb.ZONE, m.Option(z.Short(m)), m.Option(mdb.ID), arg)
}
func (z Zone) Export(m *Message, arg ...string) {
	z.Data.Export(m, mdb.ZONE, arg)
}
func (z Zone) Import(m *Message, arg ...string) {
	z.Data.Import(m, mdb.ZONE, arg)
}
func (z Zone) List(m *Message, arg ...string) *Message {
	m.Fields(len(arg), kit.Join([]string{mdb.TIME, z.Short(m), mdb.COUNT}), z.Field(m))
	z.Data.Select(m, mdb.ZONE, arg)
	if len(arg) == 0 {
		m.PushAction(z.Remove)
	}
	return m
}
