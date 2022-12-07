package providee

import (
	"fuxi/core"
)

type Providee struct {
	core.CoreService
	ProvideeServiceConfProp
}

func NewProvidee(pvid int16, name string) *Providee {
	p := &Providee{}
	p.pvid = pvid
	p.SetName(name)
	return p
}

func (self *Providee) Start() {
	//todo 从etcd拿所有的provider
	port := core.NewConnector("127.0.0.1", 8088)
	core.ServiceAddPort(self, port)
	self.CoreService.Start()
}