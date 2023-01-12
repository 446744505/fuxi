package internal

import (
	"fuxi/msg"
	"fuxi/providee"
)

type mapEventHandler struct {
	providee.ProvideeEventHandler
}

func (self *mapEventHandler) Init() {
	self.ProvideeEventHandler.Init()
	self.RegisterSvr()
	self.RegisterClient()
}

func (self *mapEventHandler) RegisterSvr() {
	self.RegisterMsg(&msg.GEnterMap{}, self.OnGEnterMap)
	self.RegisterMsg(&msg.GExitMap{}, self.OnGExitMap)
}

func (self *mapEventHandler) RegisterClient() {
	self.RegisterMsg(&msg.SEnterMap{}, nil)
	self.RegisterMsg(&msg.CGetInfo{}, self.OnCGetInfo)
}
