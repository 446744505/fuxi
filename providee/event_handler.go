package providee

import (
	"fuxi/core"
	msg "fuxi/gen"
)

type ProvideeEventHandler struct {
	core.CoreEventHandler
}

func (self *ProvideeEventHandler) Init() {
	self.RegisterMsg(&msg.PDispatch{}, nil)
	self.RegisterMsg(&msg.BindPvid{}, nil)
	self.RegisterMsg(&msg.UnBindPvid{}, nil)
}

func (self *ProvideeEventHandler) OnSessionAdd(session core.Session) {
	bind := &msg.BindPvid{}
	service := session.Port().Service()
	conf := service.(ProvideeServiceConf)
	bind.PVID = conf.PVID()
	bind.Name = conf.Name()
	session.Send(bind)
}