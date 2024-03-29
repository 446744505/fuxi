package core

import (
	"github.com/davyxu/cellnet/peer"
	"os"
	"os/signal"
	"syscall"
)

type NetControl interface {
	Controler
	Wait()
}

type Controler interface {
	Start()
	Stop()
}

type NetControlImpl struct {
	CoreNet
	peer.CoreRunningTag
	signal chan os.Signal
}

func (self *NetControlImpl) Start() {
	self.WaitStopFinished()

	if self.IsRunning() {
		return
	}

	self.CoreNet.Start()
	self.SetRunning(true)
}

func (self *NetControlImpl) Stop() {
	if !self.IsRunning() {
		return
	}

	if self.IsStopping() {
		return
	}

	self.StartStopping()
	self.CoreNet.Stop()
	self.EndStopping()
}

func (self *NetControlImpl) Wait() {
	self.signal = make(chan os.Signal)
	signal.Notify(self.signal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-self.signal
	close(self.signal)
	self.Stop()
}
