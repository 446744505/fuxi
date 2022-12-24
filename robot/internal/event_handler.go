package internal

import (
	"fuxi/core"
	"fuxi/msg"
	"github.com/davyxu/cellnet"
)

type robotEventHandler struct {
	core.CoreEventHandler
}

func (self *robotEventHandler) Init() {
	self.RegisterMsg(&msg.CEnterGame{}, nil)
	self.RegisterMsg(&msg.SEnterGame{}, self.OnSEnterGame)
	self.RegisterMsg(&msg.SEnterMap{}, self.OnSEnterMap)
}

func (self *robotEventHandler) OnSessionAdd(session core.Session) {
	ctx := session.Port().Peer().(cellnet.ContextSet)
	if role, ok := ctx.GetContext(CtxRole); ok {
		role.(*Role).OnAddSession(session)
	}
}

func (self *robotEventHandler) OnSessionRemoved(session core.Session) {
	ctx := session.Port().Peer().(cellnet.ContextSet)
	if role, ok := ctx.GetContext(CtxRole); ok {
		role.(*Role).OnRemoveSession(session)
	}
}


