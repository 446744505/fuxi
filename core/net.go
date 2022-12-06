package core

import (
	"github.com/davyxu/cellnet/peer"
	"os"
	"os/signal"
)

type ServiceCreateFunc func() Service

type Net interface {
	Start()
	Stop()
	Wait()
	Exit()
	NewService(conf ServiceConf) Service
}

type CoreNet struct {
	peer.CoreRunningTag
	services []Service
	signal chan os.Signal
}

func (self *CoreNet) Start() {
	self.WaitStopFinished()

	if self.IsRunning() {
		return
	}

	for _, service := range self.services {
		service.Start()
	}
	
	self.SetRunning(true)
}

func (self *CoreNet) Stop() {
	if !self.IsRunning() {
		return
	}
	if self.IsStopping() {
		return
	}
	self.StartStopping()
	self.WaitStopFinished()
}

func (self *CoreNet) NewService(creater ServiceCreateFunc) Service {
	s := creater()
	self.services = append(self.services, s)
	return s
}

func (self *CoreNet) Wait() {
	self.signal = make(chan os.Signal)
	signal.Notify(self.signal, os.Interrupt, os.Kill)
	<- self.signal
	self.Stop()
}