package ice

import (
	ice "github.com/shylinux/icebergs"
	_ "github.com/shylinux/icebergs/base"
	_ "github.com/shylinux/icebergs/core"
	_ "github.com/shylinux/icebergs/misc"
	kit "github.com/shylinux/toolkits"
)

func Run(arg ...string) string {
	return ice.Run(arg...)
}
func RunServe(port string, arg ...string) string {
	return ice.Run(kit.Simple("serve", "start", "dev", "", "port", port, arg)...)
}
