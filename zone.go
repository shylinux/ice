package ice

import (
	ice "shylinux.com/x/icebergs"
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
	if m.OptionFields() == "" {
		m.OptionFields(m.Config(mdb.SHORT), m.Config(mdb.FIELD))
	}
	m.Option(ice.CACHE_LIMIT, "-1")
	z.Data.Export(m, mdb.ZONE)
}
func (z Zone) Import(m *Message, arg ...string) {
	z.Data.Import(m, mdb.ZONE)
}
func (z Zone) Prev(m *Message, arg ...string) {
	mdb.PrevPage(m.Message, arg[0], arg[1:]...)
}
func (z Zone) Next(m *Message, arg ...string) {
	mdb.NextPageLimit(m.Message, arg[0], arg[1:]...)
}
func (z Zone) List(m *Message, arg ...string) *Message {
	m.Fields(len(arg), kit.Join([]string{mdb.TIME, z.Short(m), mdb.COUNT}), z.Field(m))
	if z.Data.Select(m, mdb.ZONE, arg); len(arg) == 0 {
		m.PushAction(z.Remove)
		m.StatusTimeCount()
	} else {
		m.Richs(m.PrefixKey(), "", arg[0], func(key string, value map[string]interface{}) {
			m.StatusTimeCountTotal(kit.Value(value, kit.Keym(mdb.COUNT)))
		})
	}
	return m
}
