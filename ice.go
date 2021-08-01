package ice

import (
	ice "github.com/shylinux/icebergs"
	_ "github.com/shylinux/icebergs/base"
	_ "github.com/shylinux/icebergs/core"
	_ "github.com/shylinux/icebergs/misc"
	kit "github.com/shylinux/toolkits"
)

type Message struct {
	*ice.Message
}

func Run(arg ...string) string {
	ice.Pulse.Set("detail")
	ice.Pulse.Set("append")
	ice.Pulse.Set("result")
	return ice.Run(arg...)
}
func RunServe(port string, arg ...string) string {
	return ice.Run(kit.Simple("serve", "start", "dev", "", "port", port, arg)...)
}
func RunPage(port string, path string, cb func(*Page)) {
	App("web", path, cb)
	RunServe(port)
}
