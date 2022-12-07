package internal

import (
	"fuxi/providee"
)

type mapEventHandler struct {
	providee.ProvideeEventHandler
}

func (self *mapEventHandler) Init() {
	self.ProvideeEventHandler.Init()
}
