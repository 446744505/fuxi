package provider

import (
	"fuxi/core"
	"github.com/davyxu/golog"
)

var Log = golog.New("provider")
var Provider *provider

type provider struct {
	core.CoreService

	providees map[int16]core.Session
}

func NewProvider() *provider {
	Provider = &provider{}
	Provider.SetName("provider")
	Provider.NewPort(func() core.Port {
		return core.NewAcceptor("127.0.0.1", 8088)
	})
	Provider.NewEventHandler(func() core.EventHandler {
		return &ProviderEventHandler{}
	})
	return Provider
}

func (self *provider) BindProvidee(pvid int16, session core.Session) {
	if self.providees == nil {
		self.providees = make(map[int16]core.Session)
	}
	Log.Infof("bind providee %d [%s]", pvid, session)
	self.providees[pvid] = session
}
