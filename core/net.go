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
}

type CoreNet struct {
	services []Service
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
	self.services = append(self.services, service)
}

func init() {
	proc.RegisterProcessor("fxtcp.ltv", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {
		bundle.SetTransmitter(new(CoreMessageTransmitter))
		bundle.SetHooker(new(tcp.MsgHooker))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}