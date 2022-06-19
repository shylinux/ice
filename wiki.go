package ice

import kit "shylinux.com/x/toolkits"

func WikiCmd(obj Any, arg ...Any) string {
	return Cmd(kit.Keys("web.wiki", kit.FileName(2)), obj, arg...)
}
