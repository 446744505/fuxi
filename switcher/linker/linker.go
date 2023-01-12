package linker

import (
	"fuxi/core"
	"fuxi/switcher/util"
	"github.com/davyxu/golog"
	"strconv"
	"strings"
	"sync"
)

var Log = golog.New("linker")
var Linker *linker

type linker struct {
	core.CoreService

	clients sync.Map
}

func NewLinker() *linker {
	Linker = &linker{}
	Linker.SetName("linker")
	Linker.SetEventHandler(&linkerEventHandler{})
	Linker.SetDispatcherHandler(OnDispatch)

	url := core.Args.Get("linker")
	arr := strings.Split(url, ":")
	host := arr[0]
	var port, _ = strconv.Atoi(arr[1])
	if core.ServiceAddPort(Linker, core.NewAcceptor("linker", host, port)) {
		meta := &core.SwitcherMeta{
			NodeName:    core.NodeNameLinker,
			LinkerUrl:   url,
			ProviderUrl: core.Args.Get("provider"),
		}
		core.ETCD.Put(meta.Path(), url)
	}
	return Linker
}

func (self *linker) AddSession(session core.Session) {
	self.clients.Store(session.ID(), session)
}

func (self *linker) RemoveSession(session core.Session) {
	sid := session.ID()
	self.clients.Delete(sid)
	util.ClientBroken(sid)
}

func (self *linker) GetClient(sid int64) core.Session {
	if sess, ok := self.clients.Load(sid); ok {
		return sess.(core.Session)
	}
	return nil
}

func init() {
	util.DispatchToClient = func(dispatch *core.Dispatch) {
		sid := dispatch.FTId()
		cli := Linker.GetClient(sid)
		if cli == nil {
			Log.Errorf("not exist client %d, msgId: %d session: %s",
				sid, dispatch.MsgId, dispatch.Session())
			return
		}
		cli.SendRaw(dispatch.MsgId, dispatch.MsgData)
	}
}
