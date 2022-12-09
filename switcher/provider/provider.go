package provider

import (
	"fuxi/core"
	"github.com/davyxu/golog"
	"sync"
)

var Log = golog.New("provider")
var Provider *provider

type provider struct {
	core.CoreService

	lock sync.RWMutex
	providees map[int32]core.Session
}

func NewProvider() *provider {
	Provider = &provider{}
	Provider.SetName("provider")
	Provider.SetEventHandler(&ProviderEventHandler{})
	Provider.SetDispatcherHandler(OnDispatch)
	core.ServiceAddPort(Provider, core.NewAcceptor("127.0.0.1", 8088))
	return Provider
}

func (self *provider) BindProvidee(pvid int32, name string, session core.Session) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.providees == nil {
		self.providees = make(map[int32]core.Session)
	}
	Log.Infof("bind providee [%d] [%s] [%s]", pvid, name, session)
	self.providees[pvid] = session
}

func (self *provider) UnBindProvidee(pvid int32) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if pvid > 0 && self.providees != nil {
		Log.Infof("unbind providee [%d]", pvid)
		delete(self.providees, pvid)
	}
}

func (self *provider) RemoveSession(session core.Session) {
	var deleteId int32
	for pvid, tmp := range self.providees {
		if tmp.ID() == session.ID() {
			deleteId = pvid
			break
		}
	}
	self.UnBindProvidee(deleteId)
}
