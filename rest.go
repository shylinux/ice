package ice

import (
	ice "shylinux.com/x/icebergs"
	"shylinux.com/x/icebergs/base/web"
	kit "shylinux.com/x/toolkits"
)

type Rest struct{}

func (r Rest) Get(m *Message, url string, arg ...Any) Any {
	return kit.UnMarshal(m.Cmdx(web.SPIDE, ice.DEV, web.SPIDE_RAW, web.SPIDE_GET, url, arg))
}
