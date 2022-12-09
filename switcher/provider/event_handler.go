package provider

import (
	"fuxi/core"
	"fuxi/msg"
)

type providerEventHandler struct {
	core.CoreEventHandler
}

func (self *providerEventHandler) Init() {
	self.RegisterMsg(&msg.BindPvid{}, self.OnBindPvid)
	self.RegisterMsg(&msg.UnBindPvid{}, self.OnUnBindPvid)
}

func (self *providerEventHandler) OnSessionRemoved(session core.Session) {
	Provider.RemoveSession(session)
}

