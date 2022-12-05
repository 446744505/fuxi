package net

import "github.com/davyxu/cellnet"

type NetFactory interface {
	NewNet(conf *NetConf) Neter
	NewService(conf *ServiceConf) Servicer
	NewPort(conf *PortConf) Porter
}

type ServiceFactory interface {
	NewEventHandler() EventHandler
	NewSession(raw cellnet.Session) *Session
}

type SampleNetFactory struct {
}

type SampleServiceFactory struct {
	
}

func NewSampleNetFactory() NetFactory {
	return &SampleNetFactory{}
}

func NewSampleServiceFactory() ServiceFactory {
	return &SampleServiceFactory{}
}

func (self *SampleNetFactory) NewNet(conf *NetConf) Neter {
	return NewNet(conf)
}

func (self *SampleNetFactory) NewService(conf *ServiceConf) Servicer {
	return NewService(conf)
}

func (self *SampleNetFactory) NewPort(conf *PortConf) Porter {
	if conf.Typ == Acceptor {
		return NewAccept(conf)
	} else if conf.Typ == Connector {
		return NewConnect(conf)
	}
	return nil
}

func (self *SampleServiceFactory) NewEventHandler() EventHandler {
	return &SampleEventHandler{}
}

func (self *SampleServiceFactory) NewSession(raw cellnet.Session) *Session {
	return NewSession(raw)
}
