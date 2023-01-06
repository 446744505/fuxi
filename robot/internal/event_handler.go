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
	self.RegisterMsg(&msg.CGetInfo{}, nil)
}

func (self *robotEventHandler) OnSessionAdd(session core.Session) {
	ctx := session.Port().Peer().(cellnet.ContextSet)
	if role, ok := ctx.GetContext(CtxTypeRole); ok {
		role.(*Role).OnAddSession(session)
	}
}

func (self *robotEventHandler) OnSessionRemoved(session core.Session) {
	ctx := session.Port().Peer().(cellnet.ContextSet)
	if role, ok := ctx.GetContext(CtxTypeRole); ok {
		role.(*Role).OnRemoveSession(session)
	}
}


