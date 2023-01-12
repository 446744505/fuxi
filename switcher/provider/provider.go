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

	providees sync.Map
}

func NewProvider() *provider {
	Provider = &provider{}
	Provider.SetName("provider")
	Provider.SetEventHandler(&providerEventHandler{})
	Provider.SetDispatcherHandler(OnDispatch)

	url := core.Args.Get("provider")
	arr := strings.Split(url, ":")
	host := arr[0]
	var port, _ = strconv.Atoi(arr[1])
	if core.ServiceAddPort(Provider, core.NewAcceptor("provider", host, port)) {
		meta := &core.SwitcherMeta{
			NodeName:    core.NodeNameProvider,
			LinkerUrl:   core.Args.Get("linker"),
			ProviderUrl: url,
		}
		core.ETCD.Put(meta.Path(), url)
	}
	return Provider
}

func (self *provider) BindProvidee(pvid int32, name string, session core.Session) {
	self.providees.Store(pvid, session)

	info, _ := session.GetContext(util.CtxTypeSessionInfo)
	provideeInfo := info.(*util.ProvideeSessionInfo)
	provideeInfo.Pvid = pvid
	provideeInfo.Name = name
	Log.Infof("bind providee [%d] [%s] [%s]", pvid, name, session)

	url := session.Port().HostPortString()
	meta := &core.ProvideeMeta{
		NodeName:    core.NodeNameProvidee,
		ProviderUrl: url,
		Pvid:        pvid,
	}
	core.ETCD.Put(meta.Path(), name)
}

func (self *provider) UnBindProvidee(pvid int32) {
	var url = ""
	if pvid > 0 {
		if value, ok := self.providees.Load(pvid); ok {
			self.providees.Delete(pvid)
			session := value.(core.Session)
			url = session.Port().HostPortString()
			Log.Infof("unbind providee [%d]", pvid)
			//TODO 踢客户端下线
		}
	}

	if url != "" {
		core.ETCD.Delete(fmt.Sprintf("providee/%v/%v", url, pvid))
	}
}

func (self *provider) RemoveSession(session core.Session) {
	var deleteId int32
	self.providees.Range(func(key, value interface{}) bool {
		pvid := key.(int32)
		tmp := value.(core.Session)
		if tmp.ID() == session.ID() {
			deleteId = pvid
			return false
		}
		return true
	})
	self.UnBindProvidee(deleteId)
}

func (self *provider) GetProvidee(pvid int32) core.Session {
	if sess, ok := self.providees.Load(pvid); ok {
		return sess.(core.Session)
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
			MsgId:   int32(dispatch.MsgId),
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
			MsgId:   int32(dispatch.MsgId),
			MsgData: dispatch.MsgData,
		}
		session := dispatch.Session()
		info, _ := session.GetContext(util.CtxTypeSessionInfo)
		if provideeInfo, ok := info.(*util.ProvideeSessionInfo); ok {
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

	util.ClientBroken = func(clientSid int64) {
		Provider.providees.Range(func(_, value interface{}) bool {
			session := value.(core.Session)
			session.Send(&msg.ClientBroken{ClientSid: clientSid})
			return true
		})
	}
}
