package provider

import (
	"fuxi/core"
	msg "fuxi/gen"
)

type ProviderEventHandler struct {
	core.CoreEventHandler
}

func (self *ProviderEventHandler) Init() {
	self.RegisterMsg(&msg.PDispatch{}, self.OnPDispatch)
	self.RegisterMsg(&msg.BindPvid{}, self.OnBindPvid)
	self.RegisterMsg(&msg.UnBindPvid{}, self.OnUnBindPvid)
}

