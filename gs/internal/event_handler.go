package internal

import "fuxi/providee"

type gsEventHandler struct {
	providee.ProvideeEventHandler
}

func (self *gsEventHandler) Init() {
	self.ProvideeEventHandler.Init()
}