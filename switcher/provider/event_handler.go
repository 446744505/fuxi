package provider

import (
	"fuxi/core"
	"fuxi/msg"
	"fuxi/switcher/util"
)

type providerEventHandler struct {
	core.CoreEventHandler
}

func (self *providerEventHandler) Init() {
	self.RegisterMsg(&msg.BindPvid{}, self.OnBindPvid)
	self.RegisterMsg(&msg.UnBindPvid{}, self.OnUnBindPvid)
	self.RegisterMsg(&msg.MessageBox{}, nil)
	self.RegisterMsg(&msg.ClientBroken{}, nil)
}

func (self *providerEventHandler) OnSessionAdd(session core.Session) {
	session.SetContext(util.CtxTypeSessionInfo, &util.ProvideeSessionInfo{})
}

func (self *providerEventHandler) OnSessionRemoved(session core.Session) {
	Provider.RemoveSession(session)
}
