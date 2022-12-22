package linker

import (
	"fuxi/core"
	"fuxi/switcher/linker/util"
	"github.com/davyxu/golog"
	"strconv"
	"strings"
	"sync"
)

var Log = golog.New("linker")
var Linker *linker

type linker struct {
	core.CoreService

	lock sync.RWMutex
	clients map[int64]core.Session
}

func NewLinker() *linker {
	Linker = &linker{
		clients: make(map[int64]core.Session),
	}
	Linker.SetName("linker")
	Linker.SetEventHandler(&linkerEventHandler{})
	Linker.SetDispatcherHandler(OnDispatch)

	url := core.Args.Get("linker")
	arr := strings.Split(url, ":")
	host := arr[0]
	var port, _ = strconv.Atoi(arr[1])
	core.ServiceAddPort(Linker, core.NewAcceptor("linker", host, port))
	core.ETCD.Put("linker/" + url, url)
	return Linker
}

func (self *linker) AddSession(session core.Session) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.clients[session.ID()] = session
}

func (self *linker) RemoveSession(session core.Session) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.clients, session.ID())
}

func (self *linker) GetClient(sid int64) core.Session {
	self.lock.RLock()
	defer self.lock.RUnlock()
	if sess, ok := self.clients[sid]; ok {
		return sess
	}
	return nil
}

func init() {
	util.DispatchToClient = func(dispatch *core.Dispatch) {
		sid := dispatch.ToID()
		cli := Linker.GetClient(sid)
		if cli == nil {
			Log.Errorf("not exist client %d, msgId: %d session: %s",
				sid, dispatch.MsgId, dispatch.Session())
			return
		}
		cli.SendRaw(dispatch.MsgId, dispatch.MsgData)
	}
}