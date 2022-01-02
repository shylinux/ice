package ice

import (
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"
)

type Data struct {
	_delete string `name:"delete" help:"删除"`
	remove  string `name:"remove" help:"删除"`
	prunes  string `name:"prunes" help:"清理"`
	export  string `name:"export" help:"导出"`
	_import string `name:"import" help:"导入"`
	next    string `name:"next" help:"下一页"`
	prev    string `name:"prev" help:"上一页"`
}

func (d Data) Short(m *Message) string { return m.Config(kit.MDB_SHORT) }
func (d Data) Field(m *Message) string { return m.Config(kit.MDB_FIELD) }

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
func (d Data) Inputs(m *Message, arg ...interface{}) {
	m.Cmdy(mdb.INPUTS, m.PrefixKey(), "", kit.Simple(arg))
}
func (d Data) Prunes(m *Message, arg ...interface{}) {
	m.OptionFields(m.Config(kit.META_FIELD))
	m.Cmdy(mdb.PRUNES, m.PrefixKey(), "", kit.Simple(arg))
}
func (d Data) Export(m *Message, arg ...interface{}) {
	m.OptionFields(m.Config(kit.META_FIELD))
	m.Cmdy(mdb.EXPORT, m.PrefixKey(), "", kit.Simple(arg))
}
func (d Data) Import(m *Message, arg ...interface{}) {
	m.Cmdy(mdb.IMPORT, m.PrefixKey(), "", kit.Simple(arg))
}
func (d Data) Next(m *Message, arg ...string) {
	mdb.NextPage(m.Message, arg[0], arg[1:]...)
}
func (d Data) Prev(m *Message, arg ...string) {
	mdb.PrevPage(m.Message, arg[0], arg[1:]...)
}
func (d Data) List(m *Message, arg ...string) {
	m.Echo("hello world")
}
