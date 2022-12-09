package core

import "github.com/davyxu/cellnet"

type DispatchHandler func(msg *Dispatch)

type Service interface {
	Controler
	ServiceBundle
	addPort(port Port)
}

type CoreService struct {
	CoreServiceConf
	CoreServiceBundle
	ports []Port
}

type ServiceBundle interface {
	SetEventHandler(handler EventHandler)
	EventHandler() EventHandler
	SetSessionCreater(creater SessionCreater)
	SessionCreater() SessionCreater
	SetDispatcherHandler(handler DispatchHandler)
	DispatchHandler() DispatchHandler
}

type CoreServiceBundle struct {
	evtHandler EventHandler
	creater    SessionCreater
	dipHandler DispatchHandler
}

func (self *CoreService) Start() {
	self.evtHandler.Init()
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

func (self *CoreServiceBundle) EventHandler() EventHandler {
	if self.evtHandler != nil {
		return self.evtHandler
	}
	return &CoreEventHandler{}
}

func (self *CoreServiceBundle) SetEventHandler(handler EventHandler) {
	self.evtHandler = handler
}

func (self *CoreServiceBundle) SetSessionCreater(creater SessionCreater) {
	self.creater = creater
}

func (self *CoreServiceBundle) SessionCreater() SessionCreater {
	if self.creater != nil {
		return self.creater
	}

	return func(raw cellnet.Session) Session {
		return &CoreSession{raw: raw}
	}
}

func (self *CoreServiceBundle) SetDispatcherHandler(handler DispatchHandler) {
	self.dipHandler = handler
}

func (self *CoreServiceBundle) DispatchHandler() DispatchHandler {
	return self.dipHandler
}

func ServiceAddPort(service Service, port Port) {
	port.SetService(service)
	service.addPort(port)
}
