package ice

import (
	ice "shylinux.com/x/icebergs"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/web"
	kit "shylinux.com/x/toolkits"
)

func Run(arg ...string) string {
	ice.Pulse.Set(ice.MSG_DETAIL)
	ice.Pulse.Set(ice.MSG_APPEND)
	ice.Pulse.Set(ice.MSG_RESULT)
	return ice.Run(arg...)
}
func RunServe(arg ...string) string {
	return ice.Run(kit.Simple(web.SERVE, cli.START, ice.DEV, "", arg)...)
}
