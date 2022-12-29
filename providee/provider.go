package providee

import (
	"fuxi/core"
	"sync"
)

type Provider struct {
	providerUrl string
	session     core.Session

	lock sync.RWMutex
	provideePvids map[int32]string // key=pvid value=server name
}

func NewProvider(url string) *Provider {
	return &Provider{
		providerUrl: url,
		provideePvids: make(map[int32]string),
	}
}

func (self *Provider) IsActive() bool {
	return self.session != nil
}

func (self *Provider) Send(msg core.Msg) bool {
	if !self.IsActive() {
		return false
	}
	self.session.Send(msg)
	return true
}

func (self *Provider) HaveProvidee(pvid int32) bool {
	self.lock.RLock()
	defer self.lock.RUnlock()
	if pvid == 0 {
		return len(self.provideePvids) > 0
	}
	if _, ok := self.provideePvids[pvid]; ok {
		return true
	}

	return false
}

func (self *Provider) AddProvidee(pvid int32, name string) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.provideePvids[pvid] = name
	core.Log.Infof("providee add other providee %v[%v]", pvid, name)
}

func (self *Provider) RemoveProvidee(pvid int32, name string) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.provideePvids, pvid)
	core.Log.Infof("providee remove other providee %v[%v]", pvid, name)
}

func (self *Provider) ForProvidees(cb func(int32, string)) {
	self.lock.RLock()
	defer self.lock.RUnlock()
	for pvid, name := range self.provideePvids {
		cb(pvid, name)
	}
}