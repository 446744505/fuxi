package internal

import (
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

func (self *gsEventHandler) RegisterSvr() {
	self.RegisterMsg(&msg.LEnterGame{}, self.OnLEnterGame)
	self.RegisterMsg(&msg.SEnterGame{}, nil)
}

func (self *gsEventHandler) RegisterClient() {
	self.RegisterMsg(&msg.MapNtf{}, self.OnMapNtf)
}