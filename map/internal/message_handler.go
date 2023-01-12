package internal

import (
	"fuxi/core"
	"fuxi/msg"
)

func (self *mapEventHandler) OnGEnterMap(p core.Msg) {
	enter := p.(*msg.GEnterMap)
	Map.OnRoleEnter(enter)
}

func (self *mapEventHandler) OnGExitMap(p core.Msg) {
	exit := p.(*msg.GExitMap)
	Map.OnRoleExit(exit.RoleId)
}

func (self *mapEventHandler) OnCGetInfo(p core.Msg) {
	Log.Infof("map received role %v cgetinfo", p.FTId())
}
