package ice

import (
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

type Data struct{}

func (d Data) Short(m *Message) string {
	return m.Conf(m.PrefixKey(), kit.Keym(kit.MDB_SHORT))
}
func (d Data) Field(m *Message) string {
	return m.Conf(m.PrefixKey(), kit.Keym(kit.MDB_FIELD))
}

func (d Data) Insert(m *Message, arg ...interface{}) {
	m.Cmdy(mdb.INSERT, m.PrefixKey(), "", kit.Simple(arg))
}
func (d Data) Modify(m *Message, arg ...interface{}) {
	m.Cmdy(mdb.MODIFY, m.PrefixKey(), "", kit.Simple(arg))
}
func (d Data) Delete(m *Message, arg ...interface{}) {
	m.Cmdy(mdb.DELETE, m.PrefixKey(), "", kit.Simple(arg))
}
func (d Data) Select(m *Message, arg ...interface{}) {
	m.Cmdy(mdb.SELECT, m.PrefixKey(), "", kit.Simple(arg))
}

func (d Data) Inputs(m *Message, arg ...interface{}) {
	m.Cmdy(mdb.INPUTS, m.PrefixKey(), "", kit.Simple(arg))
}

func (d Data) List(m *Message, arg ...string) {
	m.Echo("hello world")
}
