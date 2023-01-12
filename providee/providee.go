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
type OnClientBroken func(clientSid int64)

type providee struct {
	core.CoreService
	ProvideeServiceConfProp

	OnClientBroken

	onProvideeUpdate OnProvideeUpdate
	providerMap      sync.Map
}

type provideeWatcher struct {
}

func NewProvidee(pvid int32, name string) *providee {
	Providee = &providee{}
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

	provider = NewProvider(providerUrl)
	provider.session = session
	self.providerMap.Store(providerUrl, provider)
}

func (self *providee) OnRemoveSession(session core.Session) {
	providerUrl := session.Port().HostPortString()
	if provider, ok := self.providerMap.Load(providerUrl); ok {
		provider.(*Provider).session = nil
	}
}

func (self *providee) GetProvider(providerUrl string) *Provider {
	if provider, ok := self.providerMap.Load(providerUrl); ok {
		return provider.(*Provider)
	}
	return nil
}

func (self *providee) SetOnProvideeUpdate(cb OnProvideeUpdate) {
	self.onProvideeUpdate = cb
	self.providerMap.Range(func(key, value interface{}) bool {
		providerUrl := key.(string)
		provider := value.(*Provider)
		provider.ForProvidees(func(pvid int32, name string) {
			meta := &core.ProvideeMeta{
				ProviderUrl: providerUrl,
				NodeName:    name,
				Pvid:        pvid,
			}
			cb(false, meta)
		})
		return true
	})
}

func (self *providee) SendToProvidee(pvid int32, msg core.Msg) bool {
	msg.SetToType(core.MsgToProvidee)
	msg.SetFTId(int64(pvid))
	//TODO 考虑发送顺序，往同一个provider发送
	p := self.getOneProvider(pvid)
	if p == nil {
		core.Log.Errorf("provider %v can not used", pvid)
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
	var providers []*Provider
	self.providerMap.Range(func(_, value interface{}) bool {
		provider := value.(*Provider)
		if provider.IsActive() && provider.HaveProvidee(toPvid) {
			providers = append(providers, provider)
		}
		return true
	})
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
		Providee.providerMap.Store(meta.ProviderUrl, provider)
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
