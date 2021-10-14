package ice

import (
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

type Data struct {
	remove string `name:"delete" help:"删除"`
	prev   string `name:"prev" help:"上一页"`
	next   string `name:"next" help:"下一页"`
}

func (d Data) Short(m *Message) string { return m.Config(kit.MDB_SHORT) }
func (d Data) Field(m *Message) string { return m.Config(kit.MDB_FIELD) }

func (d Data) Inputs(m *Message, arg ...interface{}) {
	m.Cmdy(mdb.INPUTS, m.PrefixKey(), "", kit.Simple(arg))
}
func (d Data) Insert(m *Message, arg ...interface{}) {
	m.Cmdy(mdb.INSERT, m.PrefixKey(), "", kit.Simple(arg))
}
func (d Data) Delete(m *Message, arg ...interface{}) {
	m.Cmdy(mdb.DELETE, m.PrefixKey(), "", kit.Simple(arg))
}
func (d Data) Modify(m *Message, arg ...interface{}) {
	m.Cmdy(mdb.MODIFY, m.PrefixKey(), "", kit.Simple(arg))
}
func (d Data) Select(m *Message, arg ...interface{}) {
	m.Cmdy(mdb.SELECT, m.PrefixKey(), "", kit.Simple(arg))
}
func (d Data) Prev(m *Message, arg ...string) {
	mdb.PrevPage(m.Message, arg[0], arg[1:]...)
}
func (d Data) Next(m *Message, arg ...string) {
	mdb.NextPage(m.Message, arg[0], arg[1:]...)
}
func (d Data) List(m *Message, arg ...string) {
	m.Echo("hello world")
}
