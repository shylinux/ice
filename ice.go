package ice

import (
	ice "shylinux.com/x/icebergs"
	_ "shylinux.com/x/icebergs/base"
	_ "shylinux.com/x/icebergs/core"
	_ "shylinux.com/x/icebergs/misc"
	kit "shylinux.com/x/toolkits"
)

type Message struct{ *ice.Message }

func Run(arg ...string) string {
	ice.Pulse.Set(ice.MSG_DETAIL)
	ice.Pulse.Set(ice.MSG_APPEND)
	ice.Pulse.Set(ice.MSG_RESULT)
	return ice.Run(arg...)
}
func RunServe(port string, arg ...string) string {
	return ice.Run(kit.Simple("serve", "start", "dev", "", "port", port, arg)...)
}
