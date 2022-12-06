package core

import (
	"github.com/davyxu/cellnet/peer"
	"sync"
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
	signal sync.WaitGroup
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
	self.signal.Add(1)
	self.signal.Wait()
}

func (self *CoreNet) Exit() {
	self.signal.Done()
}