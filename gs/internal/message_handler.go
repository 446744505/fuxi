package internal

import (
	"fuxi/core"
	"fuxi/msg"
)

func (self *gsEventHandler) OnLEnterGame(p core.Msg) {
	enter := p.(*msg.LEnterGame)
	GS.OnRoleEnter(enter)
}