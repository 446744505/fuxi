package internal

import (
	"fuxi/core"
	"fuxi/msg"
	"fuxi/providee"
)

type gsEventHandler struct {
	providee.ProvideeEventHandler
}

func (self *gsEventHandler) Init() {
	self.ProvideeEventHandler.Init()
	self.RegisterSvr()
	self.RegisterClient()
}

func (self *gsEventHandler) OnSessionRemoved(session core.Session) {
	self.ProvideeEventHandler.OnSessionRemoved(session)
	GS.OnProviderBroken(session)
}

func (self *gsEventHandler) RegisterSvr() {
	self.RegisterMsg(&msg.LEnterGame{}, self.OnLEnterGame)
	self.RegisterMsg(&msg.GEnterMap{}, nil)
	self.RegisterMsg(&msg.GExitMap{}, nil)
}

func (self *gsEventHandler) RegisterClient() {
	self.RegisterMsg(&msg.SEnterGame{}, nil)
	self.RegisterMsg(&msg.CGetInfo{}, self.OnCGetInfo)
}
