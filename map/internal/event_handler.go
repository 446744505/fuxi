package internal

import "fuxi/switcher/provider"

type mapEventHandler struct {
	provider.ProviderEventHandler
}

func (self *mapEventHandler) Init() {
	self.ProviderEventHandler.Init()
}
