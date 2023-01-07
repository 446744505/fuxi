package internal

import (
	"fuxi/core"
	"fuxi/msg"
)

func (self *gsEventHandler) OnLEnterGame(p core.Msg) {
	enter := p.(*msg.LEnterGame)
	GS.OnRoleEnter(enter)
}


func (self *gsEventHandler) OnCGetInfo(p core.Msg) {
	Log.Infof("gs received role %v cgetinfo", p.FTId())
}