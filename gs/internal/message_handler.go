package internal

import (
	"fuxi/core"
	"fuxi/msg"
)

func (self *gsEventHandler) OnMapNtf(p core.Msg) {
	ntf := p.(*msg.MapNtf)
	Log.Infof("receive map notify, pvid:%d", ntf.PVID)
}

func (self *gsEventHandler) OnLEnterGame(p core.Msg) {
	enter := p.(*msg.LEnterGame)
	GS.OnRoleEnter(enter)
}