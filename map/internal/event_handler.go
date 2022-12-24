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
	self.RegisterMsg(&msg.GEnterMap{}, self.OnGEnterMap)
	self.RegisterMsg(&msg.SEnterMap{}, nil)
}

func (self *mapEventHandler) OnSessionRemoved(session core.Session) {
	self.ProvideeEventHandler.OnSessionRemoved(session)
	//todo 踢当前session的所有role下线
}
