package providee

import (
	"fuxi/core"
	"math/rand"
	"sync"
)

var Providee *providee

type providee struct {
	core.CoreService
	ProvideeServiceConfProp

	lock sync.RWMutex
	providerMap map[int64]core.Session
	providerList []core.Session
}

func NewProvidee(pvid int32, name string) *providee {
	Providee = &providee{
		providerMap: make(map[int64]core.Session),
	}
	Providee.pvid = pvid
	Providee.SetName(name)
	Providee.SetEventHandler(&ProvideeEventHandler{})
	return Providee
}

func (self *providee) Start() {
	//todo 从etcd拿所有的provider
	port := core.NewConnector("127.0.0.1", 8088)
	core.ServiceAddPort(self, port)
	self.CoreService.Start()
}

func (self *providee) OnAddSession(session core.Session) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.providerMap[session.ID()] = session
	self.providerList = append(self.providerList, session)
}

func (self *providee) OnRemoveSession(session core.Session) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.providerMap, session.ID())
	for i := 0; i < len(self.providerList); i++ {
		if self.providerList[i].ID() == session.ID() {
			self.providerList = append(self.providerList[:i], self.providerList[i+1:]...)
			break
		}
	}
}

//todo 先随便选一个provider发往providee，后续要考虑顺序，其他providee已与某个provider断开
func (self *providee) SendToProvidee(pvid int32, msg core.Msg) {
	msg.SetToType(core.MsgToProvidee)
	msg.SetToID(int64(pvid))
	p := self.getOneProvider()
	if p == nil {
		core.Log.Errorln("not any provider can be used")
		return
	}
	p.Send(msg)
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
