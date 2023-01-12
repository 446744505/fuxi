package providee

import (
	"fuxi/core"
	"sync"
)

type Provider struct {
	providerUrl string
	session     core.Session

	providees sync.Map // key=pvid value=server name
}

func NewProvider(url string) *Provider {
	return &Provider{
		providerUrl: url,
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
	if pvid == 0 {
		have := false
		self.providees.Range(func(key, value interface{}) bool {
			have = true
			return false //有任意一个就跳出循环
		})
		return have
	}

	if _, ok := self.providees.Load(pvid); ok {
		return true
	}

	return false
}

func (self *Provider) AddProvidee(pvid int32, name string) {
	self.providees.Store(pvid, name)
	core.Log.Infof("providee add other providee %v[%v]", pvid, name)
}

func (self *Provider) RemoveProvidee(pvid int32, name string) {
	self.providees.Delete(pvid)
	core.Log.Infof("providee remove other providee %v[%v]", pvid, name)
}

func (self *Provider) ForProvidees(cb func(int32, string)) {
	self.providees.Range(func(k, v interface{}) bool {
		pvid := k.(int32)
		name := v.(string)
		cb(pvid, name)
		return true
	})
}
