package linker

import (
	"fuxi/core"
	"fuxi/msg"
	"fuxi/switcher/util"
)

type linkerEventHandler struct {
	core.CoreEventHandler
}

func (self *linkerEventHandler) Init() {
	self.RegisterMsg(&msg.CEnterGame{}, self.OnCEnterGame)
	self.RegisterMsg(&msg.LEnterGame{}, nil)
}

func (self *linkerEventHandler) OnSessionAdd(session core.Session) {
	session.SetContext(util.CtxTypeSessionInfo, &util.LinkerSessionInfo{})
	Linker.AddSession(session)
}

func (self *linkerEventHandler) OnSessionRemoved(session core.Session) {
	Linker.RemoveSession(session)
}
