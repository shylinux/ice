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
	insert string `name:"insert zone=hi type=hello name=world text=nice" help:"添加"`
	list   string `name:"list zone id auto insert" help:"存储"`
}

func (z Zone) Inputs(m *Message, arg ...string) { mdb.ZoneInputs(m.Message, arg) }
func (z Zone) Create(m *Message, arg ...string) { mdb.ZoneCreate(m.Message, arg) }
func (z Zone) Remove(m *Message, arg ...string) { mdb.ZoneRemove(m.Message, arg) }
func (z Zone) Insert(m *Message, arg ...string) { mdb.ZoneInsert(m.Message, arg) }
func (z Zone) Modify(m *Message, arg ...string) { mdb.ZoneModify(m.Message, arg) }
func (z Zone) Export(m *Message, arg ...string) { mdb.ZoneExport(m.Message, arg) }
func (z Zone) Import(m *Message, arg ...string) { mdb.ZoneImport(m.Message, arg) }
func (z Zone) List(m *Message, arg ...string) *Message {
	mdb.ZoneSelect(m.Message, arg...)
	return m
}
func (z Zone) ListPage(m *Message, arg ...string) *Message {
	mdb.ZoneSelectPage(m.Message, arg...).Action(kit.Select("", mdb.PAGE, len(kit.Slice(arg, 0, 2)) > 0))
	return m
}
func (z Zone) Next(m *Message, arg ...string) {
	mdb.NextPageLimit(m.Message, kit.Select("0", arg, 0), kit.Slice(arg, 1)...)
}
