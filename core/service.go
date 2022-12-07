package core

type Service interface {
	Controler
	addPort(port Port)
	SetEventHandler(handler EventHandler)
	EventHandler() EventHandler
}

type CoreService struct {
	CoreServiceConf
	ports []Port
	handler EventHandler
}

func (self *CoreService) Start() {
	self.handler.Init()
	for _, port := range self.ports {
		port.Start()
	}
}

func (self *CoreService) Stop() {
	for _, port := range self.ports {
		port.Stop()
	}
}

func (self *CoreService) addPort(port Port) {
	self.ports = append(self.ports, port)
}

func (self *CoreService) EventHandler() EventHandler {
	return self.handler
}

func (self *CoreService) SetEventHandler(handler EventHandler) {
	self.handler = handler
}

func ServiceAddPort(service Service, port Port) {
	port.SetService(service)
	service.addPort(port)
}
