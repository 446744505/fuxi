package internal

import (
	"fuxi/core"
	"fuxi/msg"
)

func (self *mapEventHandler) OnGEnterMap(p core.Msg) {
	enter := p.(*msg.GEnterMap)
	Log.Infof("role %v enter map", enter.RoleId)
}
