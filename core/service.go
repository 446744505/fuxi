package core

import (
	"github.com/davyxu/cellnet"
	"sync"
)

type DispatchHandler func(msg *Dispatch)

type Service interface {
	Controler
	ServiceBundle
	Name() string
	addPort(port Port)
}

type CoreService struct {
	CoreServiceConf
	CoreServiceBundle

	lock sync.RWMutex
	ports map[string]Port
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
	self.lock.RLock()
	for _, port := range self.ports {
		port.Start()
	}
	self.lock.RUnlock()
}

func (self *CoreService) Stop() {
	self.lock.RLock()
	for _, port := range self.ports {
		port.Stop()
	}
	self.lock.RUnlock()
}

func (self *CoreService) addPort(port Port) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.ports == nil {
		self.ports = make(map[string]Port)
	}
	if old, ok := self.ports[port.Name()]; ok {
		old.Stop()
	}
	self.ports[port.Name()] = port
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
