package ice

import (
	"github.com/shylinux/icebergs/base/mdb"
	kit "github.com/shylinux/toolkits"
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
