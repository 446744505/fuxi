package internal

import (
	"fuxi/core"
	"fuxi/msg"
)

type robotEventHandler struct {
	core.CoreEventHandler
}

func (self *robotEventHandler) Init() {
	self.RegisterMsg(&msg.CEnterGame{}, nil)
	self.RegisterMsg(&msg.SEnterGame{}, self.OnSEnterGame)
}

func (self *robotEventHandler) OnSessionAdd(session core.Session) {
	Robot.AddLinkerSession(session)
}

func (self *robotEventHandler) OnSessionRemoved(session core.Session) {
	Robot.RemoveLinkerSession(session)
}


