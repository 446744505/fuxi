package providee

import "fuxi/net"

type ProvideeFactory struct {
	net.SampleNetFactory
	net.SampleServiceFactory
}

func (self *ProvideeFactory) NewEventHandler() net.EventHandler {
	return &ProvideeEventHandler{}
}
