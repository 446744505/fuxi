package provider

import (
	"fuxi/core"
	"fuxi/msg"
)

type ProviderEventHandler struct {
	core.CoreEventHandler
}

func (self *ProviderEventHandler) Init() {
	self.RegisterMsg(&msg.PDispatch{}, self.OnPDispatch)
	self.RegisterMsg(&msg.BindPvid{}, self.OnBindPvid)
	self.RegisterMsg(&msg.UnBindPvid{}, self.OnUnBindPvid)
}

func (self *ProviderEventHandler) OnSessionRemoved(session core.Session) {
	Provider.RemoveSession(session)
}

