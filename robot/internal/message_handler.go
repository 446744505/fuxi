package internal

import (
	"fuxi/core"
	"fuxi/msg"
)

func (self *robotEventHandler) OnSEnterGame(p core.Msg) {
	enter := p.(*msg.SEnterGame)
	Log.Infof("role %s enter client game", enter.Name)
}
