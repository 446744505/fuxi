package providee

import (
	"fuxi/core"
)

var Providee *providee

type providee struct {
	core.CoreService
	ProvideeServiceConfProp

	providers map[int64]core.Session
}

func NewProvidee(pvid int32, name string) *providee {
	Providee = &providee{
		providers: make(map[int64]core.Session),
	}
	Providee.pvid = pvid
	Providee.SetName(name)
	return Providee
}

func (self *providee) Start() {
	//todo 从etcd拿所有的provider
	port := core.NewConnector("127.0.0.1", 8088)
	core.ServiceAddPort(self, port)
	self.CoreService.Start()
}

func (self *providee) OnAddSession(session core.Session) {
	self.providers[session.ID()] = session
}

func (self *providee) BroadcastToProvider(msg core.Msg) {
	for _, session := range self.providers {
		session.Send(msg)
	}
}