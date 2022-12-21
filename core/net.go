package core

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
	"github.com/davyxu/golog"
)

var Log = golog.New("net")

type Net interface {
	Controler
	AddService(service Service)
	GetService(name string) Service
}

type CoreNet struct {
	services map[string]Service
}

func (self *CoreNet) Start() {
	for _, service := range self.services {
		service.Start()
	}
}

func (self *CoreNet) Stop() {
	for _, service := range self.services {
		service.Stop()
	}
}

func (self *CoreNet) AddService(service Service) {
	if self.services == nil {
		self.services = make(map[string]Service)
	}
	self.services[service.Name()] = service
}

func (self *CoreNet) GetService(name string) Service {
	if svr, ok := self.services[name]; ok {
		return svr
	}
	return nil
}

func init() {
	RegisterMsg(&Dispatch{})
	proc.RegisterProcessor("fxtcp.ltv", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {
		bundle.SetTransmitter(new(CoreMessageTransmitter))
		bundle.SetHooker(new(tcp.MsgHooker))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}