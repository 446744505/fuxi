package provider

import (
	"fmt"
	"fuxi/core"
	"fuxi/msg"
	"fuxi/switcher/util"
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

	url := core.Args.Get("provider")
	arr := strings.Split(url, ":")
	host := arr[0]
	var port, _ = strconv.Atoi(arr[1])
	if core.ServiceAddPort(Provider, core.NewAcceptor("provider", host, port)) {
		meta := &core.SwitcherMeta{
			NodeName: core.NodeNameProvider,
			LinkerUrl: core.Args.Get("linker"),
			ProviderUrl: url,
		}
		core.ETCD.Put(meta.Path(), url)
	}
	return Provider
}

func (self *provider) BindProvidee(pvid int32, name string, session core.Session) {
	self.lock.Lock()
	self.providees[pvid] = session
	self.lock.Unlock()

	info, _ := session.GetContext(util.CtxTypeSessionInfo)
	provideeInfo := info.(*util.ProvideeSessionInfo)
	provideeInfo.Pvid = pvid
	Log.Infof("bind providee [%d] [%s] [%s]", pvid, name, session)

	url := session.Port().HostPortString()
	meta := &core.ProvideeMeta{
		NodeName: core.NodeNameProvidee,
		ProviderUrl: url,
		Pvid: pvid,
	}
	core.ETCD.Put(meta.Path(), name)
}

func (self *provider) UnBindProvidee(pvid int32) {
	var url = ""
	self.lock.Lock()
	if pvid > 0 && self.providees != nil {
		if session, ok := self.providees[pvid]; ok {
			delete(self.providees, pvid)
			url = session.Port().HostPortString()
			Log.Infof("unbind providee [%d]", pvid)
		}
	}
	self.lock.Unlock()
	if url != "" {
		core.ETCD.Delete(fmt.Sprintf("providee/%v/%v", url, pvid))
	}
}

func (self *provider) RemoveSession(session core.Session) {
	var deleteId int32
	self.lock.RLock()
	for pvid, tmp := range self.providees {
		if tmp.ID() == session.ID() {
			deleteId = pvid
			break
		}
	}
	self.lock.RUnlock()
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
	util.ClientToProvidee = func(dispatch *core.Dispatch) {
		toPVID := int32(dispatch.FTId())
		prov := Provider.GetProvidee(toPVID)
		if prov == nil {
			Log.Errorf("ClientToProvidee not exist providee %d, msgId: %d session: %s",
				toPVID, dispatch.MsgId, dispatch.Session())
			return
		}

		to := &msg.MessageBox{
			MsgId: int32(dispatch.MsgId),
			MsgData: dispatch.MsgData,
		}
		session := dispatch.Session()
		info, _ := session.GetContext(util.CtxTypeSessionInfo)
		if linkerInfo, ok := info.(*util.LinkerSessionInfo); ok {
			to.UniqId = linkerInfo.RoleId
		}
		prov.Send(to)
	}

	util.ProvideeToProvidee = func(dispatch *core.Dispatch) {
		toPVID := int32(dispatch.FTId())
		prov := Provider.GetProvidee(toPVID)
		if prov == nil {
			Log.Errorf("ProvideeToProvidee not exist providee %d, msgId: %d session: %s",
				toPVID, dispatch.MsgId, dispatch.Session())
			return
		}

		to := &msg.MessageBox{
			MsgId: int32(dispatch.MsgId),
			MsgData: dispatch.MsgData,
		}
		session := dispatch.Session()
		info, _ := session.GetContext(util.CtxTypeSessionInfo)
		if provideeInfo, ok := info.(*util.ProvideeSessionInfo); ok{
			to.UniqId = int64(provideeInfo.Pvid)
		}
		prov.Send(to)
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
