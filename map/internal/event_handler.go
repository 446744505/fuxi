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
	self.RegisterMsg(&msg.MapNtf{}, nil)
}

func (self *mapEventHandler) OnSessionAdd(session core.Session) {
	self.ProvideeEventHandler.OnSessionAdd(session)
	service := session.Port().Service()
	conf := service.(providee.ProvideeServiceConf)
	ntf := &msg.MapNtf{}
	ntf.PVID = conf.PVID()
	session.Send(ntf)
}
