package linker

import (
	"fuxi/core"
	msg "fuxi/gen"
)

type LinkerEventHandler struct {
	core.CoreEventHandler
}

func (self *LinkerEventHandler) Init() {
	self.RegisterMsg(&msg.Dispatch{}, self.OnDispatch)
}
