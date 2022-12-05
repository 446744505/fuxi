package net

import (
	"sync"
)

type Neter interface {
	Start()
	Stop()
	Wait()
	StopSignal()
}

type NetOptions struct {
	factory NetFactory
	conf *NetConf
}

type Net struct {
	Services []Servicer
	Conf *NetConf
	endSignal sync.WaitGroup
}

type NetConf struct {
	Name string
	ServiceConfs []*ServiceConf

	Ops *NetOptions
}

func FXNet(ops NetOptions) Neter {
	return ops.factory.NewNet(ops.conf)
}

func NewNetConf(name string) *NetConf {
	return &NetConf{Name: name}
}

func NewNet(conf *NetConf) Neter {
	net := &Net{Conf: conf}
	net.Init()
	return net
}

func DefaultOptions(conf *NetConf) *NetOptions {
	ops := NewNetOptions(conf)
	ops.BuildFactory(NewSampleNetFactory())
	return ops
}

func NewNetOptions(conf *NetConf) *NetOptions {
	ops := &NetOptions{conf: conf}
	conf.Ops = ops
	for _, serviceConf := range conf.ServiceConfs {
		serviceConf.Ops = ops
		for _, portConf := range serviceConf.PortConfs {
			portConf.Ops = ops
		}
	}
	return ops
}

func (self *NetOptions) BuildFactory(factory NetFactory) NetOptions {
	self.factory = factory
	return *self
}

func (self *Net) Init() {
	factory := self.Conf.Ops.factory
	for _, serviceConf := range self.Conf.ServiceConfs {
		service := factory.NewService(serviceConf)
		self.Services = append(self.Services, service)
	}
}

func (self *Net) Start() {
	for _, service := range self.Services {
		service.Start()
	}
	self.endSignal.Add(1)
}

func (self *Net) Stop() {
	for _, service := range self.Services {
		service.Stop()
	}
}

func (self *Net) Wait() {
	self.endSignal.Wait()
}

func (self *Net) StopSignal() {
	self.endSignal.Done()
}

func (self *NetConf) NewService(name string) *ServiceConf {
	newConf := NewServiceConf(name)
	newConf.Factory = NewSampleServiceFactory()
	self.ServiceConfs = append(self.ServiceConfs, newConf)
	return newConf
}