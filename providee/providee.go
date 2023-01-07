package providee

import (
	"fuxi/core"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

var Providee *providee
type OnProvideeUpdate func(isRemove bool, meta *core.ProvideeMeta)

type providee struct {
	core.CoreService
	ProvideeServiceConfProp
	onProvideeUpdate OnProvideeUpdate

	lock sync.RWMutex
	providerMap map[string]*Provider
}

type provideeWatcher struct {

}

func NewProvidee(pvid int32, name string) *providee {
	Providee = &providee{
		providerMap: make(map[string]*Provider),
	}
	Providee.pvid = pvid
	Providee.SetName(name)
	Providee.SetEventHandler(&ProvideeEventHandler{})
	return Providee
}

func (self *providee) Start() {
	core.ETCD.Watch(core.NodeNameProvider, self)
	core.ETCD.Watch(core.NodeNameProvidee, &provideeWatcher{})
	self.CoreService.Start()
}

func (self *providee) OnAddSession(session core.Session) {
	providerUrl := session.Port().HostPortString()
	provider := self.GetProvider(providerUrl)
	if provider != nil {
		provider.session = session
		return
	}

	self.lock.Lock()
	defer self.lock.Unlock()
	provider = NewProvider(providerUrl)
	provider.session = session
	self.providerMap[providerUrl] = provider
}

func (self *providee) OnRemoveSession(session core.Session) {
	self.lock.Lock()
	defer self.lock.Unlock()
	providerUrl := session.Port().HostPortString()
	if provider, ok := self.providerMap[providerUrl]; ok {
		provider.session = nil
	}
}

func (self *providee) GetProvider(providerUrl string) *Provider {
	self.lock.RLock()
	defer self.lock.RUnlock()
	if provider, ok := self.providerMap[providerUrl]; ok {
		return provider
	}
	return nil
}

func (self *providee) SetOnProvideeUpdate(cb OnProvideeUpdate) {
	self.onProvideeUpdate = cb
	self.lock.RLock()
	defer self.lock.RUnlock()
	for providerUrl, provider := range self.providerMap {
		provider.ForProvidees(func(pvid int32, name string) {
			meta := &core.ProvideeMeta{
				ProviderUrl: providerUrl,
				NodeName: name,
				Pvid: pvid,
			}
			cb(false, meta)
		})
	}
}

func (self *providee) SendToProvidee(pvid int32, msg core.Msg) bool {
	msg.SetToType(core.MsgToProvidee)
	msg.SetFTId(int64(pvid))
	p := self.getOneProvider(pvid)
	if p == nil {
		core.Log.Errorln("not any provider can be used")
		return false
	}
	p.Send(msg)
	return true
}

func (self *providee) SendToProvidees(pvids []int32, msg core.Msg) {
	for _, pvid := range pvids {
		msg.NewHead()
		self.SendToProvidee(pvid, msg)
	}
}

func (self *providee) getOneProvider(toPvid int32) *Provider {
	self.lock.RLock()
	defer self.lock.RUnlock()
	var providers []*Provider
	for _, provider := range self.providerMap {
		if provider.IsActive() && provider.HaveProvidee(toPvid) {
			providers = append(providers, provider)
		}
	}
	if len(providers) == 0 {
		return nil
	}
	return providers[rand.Intn(len(providers))]
}

func (self *providee) OnAdd(key, val string) {
	arr := strings.Split(val, ":")
	host := arr[0]
	port, _ := strconv.Atoi(arr[1])
	porter := core.NewConnector(key, host, port)
	if core.ServiceAddPort(self, porter) {
		porter.Start()
	}
}

func (self *providee) OnDelete(_, _ string) {

}

func (self *provideeWatcher) OnAdd(key, val string) {
	meta := &core.ProvideeMeta{}
	meta.ValueOf(key, val)
	provider := Providee.GetProvider(meta.ProviderUrl)
	if provider == nil {
		provider = NewProvider(meta.ProviderUrl)
		Providee.lock.Lock()
		Providee.providerMap[meta.ProviderUrl] = provider
		Providee.lock.Unlock()
	}
	provider.AddProvidee(meta.Pvid, val)
	if Providee.onProvideeUpdate != nil {
		Providee.onProvideeUpdate(false, meta)
	}
}

func (self *provideeWatcher) OnDelete(key, val string) {
	meta := &core.ProvideeMeta{}
	meta.ValueOf(key, val)
	provider := Providee.GetProvider(meta.ProviderUrl)
	if provider != nil {
		provider.RemoveProvidee(meta.Pvid, meta.NodeName)
	}
	if Providee.onProvideeUpdate != nil {
		Providee.onProvideeUpdate(true, meta)
	}
}