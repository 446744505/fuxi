package providee

import (
	"fuxi/core"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

var Providee *providee

type providee struct {
	core.CoreService
	ProvideeServiceConfProp

	lock sync.RWMutex
	providerMap map[string]core.Session
	providerList []core.Session
}

func NewProvidee(pvid int32, name string) *providee {
	Providee = &providee{
		providerMap: make(map[string]core.Session),
	}
	Providee.pvid = pvid
	Providee.SetName(name)
	Providee.SetEventHandler(&ProvideeEventHandler{})
	return Providee
}

func (self *providee) Start() {
	core.ETCD.Watch(core.NodeNameProvider, self)
	self.CoreService.Start()
}

func (self *providee) OnAddSession(session core.Session) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.providerMap[session.Port().HostPortString()] = session
	self.providerList = append(self.providerList, session)
}

func (self *providee) OnRemoveSession(session core.Session) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.providerMap, session.Port().HostPortString())
	for i := 0; i < len(self.providerList); i++ {
		if self.providerList[i].Port().HostPortString() == session.Port().HostPortString() {
			self.providerList = append(self.providerList[:i], self.providerList[i+1:]...)
			break
		}
	}
}

func (self *providee) GetProvider(providerName string) core.Session {
	self.lock.RLock()
	defer self.lock.RUnlock()
	if provider, ok := self.providerMap[providerName]; ok {
		return provider
	}
	return nil
}

//todo 先随便选一个provider发往providee，后续要考虑顺序，其他providee已与某个provider断开
func (self *providee) SendToProvidee(pvid int32, msg core.Msg) bool {
	msg.SetToType(core.MsgToProvidee)
	msg.SetToID(int64(pvid))
	p := self.getOneProvider()
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

func (self *providee) getOneProvider() core.Session {
	self.lock.RLock()
	defer self.lock.RUnlock()
	i := rand.Intn(len(self.providerList))
	return self.providerList[i]
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

func (self *providee) OnDelete(key, val string) {

}