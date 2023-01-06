package provider

import (
	"fuxi/core"
	"fuxi/msg"
	"fuxi/switcher/util"
	"github.com/davyxu/cellnet"
)

type providerEventHandler struct {
	core.CoreEventHandler
}

func (self *providerEventHandler) Init() {
	self.RegisterMsg(&msg.BindPvid{}, self.OnBindPvid)
	self.RegisterMsg(&msg.UnBindPvid{}, self.OnUnBindPvid)
}

func (self *providerEventHandler) OnSessionAdd(session core.Session) {
	peer := session.Port().Peer()
	ctx := peer.(cellnet.ContextSet)
	ctx.SetContext(util.CtxTypeSessionInfo, &util.ProvideeSessionInfo{})
}

func (self *providerEventHandler) OnSessionRemoved(session core.Session) {
	Provider.RemoveSession(session)
}

