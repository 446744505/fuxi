package provider

import (
	"fuxi/core"
	"fuxi/switcher"
	"fuxi/switcher/linker/util"
	"github.com/davyxu/golog"
	"strconv"
	"strings"
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
	Provider = &provider{
		providees: make(map[int32]core.Session),
	}
	Provider.SetName("provider")
	Provider.SetEventHandler(&providerEventHandler{})
	Provider.SetDispatcherHandler(OnDispatch)

	url := switcher.Args.String("provider")
	arr := strings.Split(url, ":")
	host := arr[0]
	var port, _ = strconv.Atoi(arr[1])
	core.ServiceAddPort(Provider, core.NewAcceptor(host, port))
	core.ETCD.Put("provider/" + url, url)
	return Provider
}

func (self *provider) BindProvidee(pvid int32, name string, session core.Session) {
	self.lock.Lock()
	defer self.lock.Unlock()
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

func (self *provider) GetProvidee(pvid int32) core.Session {
	self.lock.RLock()
	defer self.lock.RUnlock()
	if sess, ok := self.providees[pvid]; ok {
		return sess
	}
	return nil
}

func init() {
	util.DispatchToProvidee = func(dispatch *core.Dispatch) {
		toPVID := int32(dispatch.ToID())
		prov := Provider.GetProvidee(toPVID)
		if prov == nil {
			Log.Errorf("dispatch not exist providee %d, msgId: %d session: %s",
				toPVID, dispatch.MsgId, dispatch.Session())
			return
		}
		prov.SendRaw(dispatch.MsgId, dispatch.MsgData)
	}
	util.SendToProvidee = func(pvid int32, msg core.Msg) {
		prov := Provider.GetProvidee(pvid)
		if prov == nil {
			Log.Errorf("send not exist providee %d, msgId: %d", pvid, msg.ID())
			return
		}
		prov.Send(msg)
	}
}
