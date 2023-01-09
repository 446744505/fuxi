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
	addPort(port Port) bool
}

type CoreService struct {
	CoreServiceConf
	CoreServiceBundle

	ports sync.Map
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
	self.ports.Range(func(_, value interface{}) bool {
		value.(Port).Start()
		return true
	})
}

func (self *CoreService) Stop() {
	self.ports.Range(func(_, value interface{}) bool {
		value.(Port).Stop()
		return true
	})
}

func (self *CoreService) addPort(port Port) bool {
	if _, ok := self.ports.Load(port.Name()); ok {
		return false
	}
	self.ports.Store(port.Name(), port)
	return true
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
		return NewCoreSession(raw)
	}
}

func (self *CoreServiceBundle) SetDispatcherHandler(handler DispatchHandler) {
	self.dipHandler = handler
}

func (self *CoreServiceBundle) DispatchHandler() DispatchHandler {
	return self.dipHandler
}

func ServiceAddPort(service Service, port Port) bool {
	port.SetService(service)
	return service.addPort(port)
}
