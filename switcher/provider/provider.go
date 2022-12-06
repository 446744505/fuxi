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

func (self *provider) BindProvidee(pvid int16, name string, session core.Session) {
	if self.providees == nil {
		self.providees = make(map[int16]core.Session)
	}
	Log.Infof("bind providee %d[%s] %s", pvid, name, session)
	self.providees[pvid] = session
}

func (self *provider) UnBindProvidee(pvid int16) {
	if pvid > 0 && self.providees != nil {
		Log.Infof("unbind providee %d", pvid)
		delete(self.providees, pvid)
	}
}

func (self *provider) RemoveSession(session core.Session) {
	var deleteId int16
	for pvid, tmp := range self.providees {
		if tmp.ID() == session.ID() {
			deleteId = pvid
			break
		}
	}
	self.UnBindProvidee(deleteId)
}
