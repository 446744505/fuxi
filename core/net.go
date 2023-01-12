package core

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
	"github.com/davyxu/golog"
	"net/http"
	_ "net/http/pprof"
)

var Log = golog.New("net")

type CoreNet struct {
	pprofPort int
	services  map[string]Service
}

func (self *CoreNet) Start() {
	if self.pprofPort > 0 {
		go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", self.pprofPort), nil)
	}
	for _, service := range self.services {
		service.Start()
	}
}

func (self *CoreNet) Stop() {
	for _, service := range self.services {
		service.Stop()
	}
}

func (self *CoreNet) SetPProfPort(port int) {
	self.pprofPort = port
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
