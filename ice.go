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
	ice.Pulse.Set(ice.MSG_DETAIL)
	ice.Pulse.Set(ice.MSG_APPEND)
	ice.Pulse.Set(ice.MSG_RESULT)
	return ice.Run(arg...)
}
func RunServe(port string, arg ...string) string {
	return ice.Run(kit.Simple("serve", "start", "dev", "", "port", port, arg)...)
}
