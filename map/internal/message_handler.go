package internal

import (
	"fuxi/core"
	"fuxi/msg"
)

func (self *mapEventHandler) OnGEnterMap(p core.Msg) {
	enter := p.(*msg.GEnterMap)
	role := &NetRole{}
	role.OnEnterMap(enter)
}
