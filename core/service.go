package core

type PortCreateFunc func() Port
type EventHandlerCreateFunc func() EventHandler

type Service interface {
	Start()
	Stop()
	NewPort(creater PortCreateFunc) Port
	NewEventHandler(creater EventHandlerCreateFunc) EventHandler
	EventHandler() EventHandler
}

type CoreService struct {
	CoreServiceConf
	ports []Port
	handler EventHandler
}

func (self *CoreService) Start() {
	for _, port := range self.ports {
		port.Start()
	}
}

func (self *CoreService) Stop() {
	for _, port := range self.ports {
		port.Stop()
	}
}

func (self *CoreService) NewPort(creater PortCreateFunc) Port {
	t := creater()
	t.SetService(self)
	self.ports = append(self.ports, t)
	return t
}

func (self *CoreService) EventHandler() EventHandler {
	return self.handler
}

func (self *CoreService) NewEventHandler(creater EventHandlerCreateFunc) EventHandler {
	self.handler = creater()
	self.handler.Init()
	return self.handler
}
