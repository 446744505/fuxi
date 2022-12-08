package internal

import (
	"fuxi/msg"
	"fuxi/providee"
)

type gsEventHandler struct {
	providee.ProvideeEventHandler
}

func (self *gsEventHandler) Init() {
	self.ProvideeEventHandler.Init()
	self.RegisterMsg(&msg.MapNtf{}, self.OnMapNtf)
}