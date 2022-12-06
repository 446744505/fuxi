package providee

import (
	"fuxi/core"
)

type Providee struct {
	core.CoreNet
}

type ProvideeService struct {
	core.CoreService
	ProvideeServiceConfProp
}

func NewProvidee(pvid int16, name string) *Providee {
	p := &Providee{}
	s := p.NewService(func() core.Service {
		s := &ProvideeService{}
		s.pvid = pvid
		s.SetName(name)
		return s
	})
	s.NewEventHandler(func() core.EventHandler {
		return &ProvideeEventHandler{}
	})
	//todo 从etcd拿所有的provider
	s.NewPort(func() core.Port {
		return core.NewConnector("127.0.0.1", 8088)
	}).SetService(s) //fixme
	return p
}