package internal

import (
	"fuxi/core"
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

func (self *mapEventHandler) OnSessionRemoved(session core.Session) {
	self.ProvideeEventHandler.OnSessionRemoved(session)
	//todo 踢当前session的所有role下线
}

func (self *mapEventHandler) RegisterSvr() {
	self.RegisterMsg(&msg.GEnterMap{}, self.OnGEnterMap)
}

func (self *mapEventHandler) RegisterClient() {
	self.RegisterMsg(&msg.SEnterMap{}, nil)
	self.RegisterMsg(&msg.CGetInfo{}, self.OnCGetInfo)
}
