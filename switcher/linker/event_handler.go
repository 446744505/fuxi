package linker

import (
	"fuxi/core"
	"fuxi/msg"
)

type LinkerEventHandler struct {
	core.CoreEventHandler
}

func (self *LinkerEventHandler) Init() {
	self.RegisterMsg(&msg.Dispatch{}, self.OnDispatch)
}
