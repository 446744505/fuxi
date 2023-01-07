package core

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
	"sync"
	"time"
)

type CtxType int

type Port interface {
	Controler
	Peer() cellnet.Peer
	SetService(service Service)
	Service() Service
	Name() string
	HostPortString() string
}

type CorePort struct {
	CorePortConf
	service Service

	pLock sync.Mutex
	peer  cellnet.Peer

	sLock sync.RWMutex
	sessions map[int64]Session
}

func (self *CorePort) Start() {
	self.pLock.Lock()
	defer self.pLock.Unlock()
	self.peer = peer.NewGenericPeer(self.typ, "", self.HostPortString(), nil)
	evtHandler := self.Service().EventHandler()
	creater := self.Service().SessionCreater()
	dipHandler := self.Service().DispatchHandler()
	proc.BindProcessorHandler(self.peer, "fxtcp.ltv", func(ev cellnet.Event) {
		switch msg := ev.Message().(type) {
		case *cellnet.SessionAccepted, *cellnet.SessionConnected:
			session := creater(ev.Session())
			session.SetPort(self)
			self.AddSession(session.ID(), session)
			evtHandler.OnSessionAdd(session)
		case *cellnet.SessionClosed:
			if session, ok := self.GetSession(ev.Session().ID()); ok {
				Log.Infof("session %s closed, reason: %s", session, msg.Reason)
				evtHandler.OnSessionRemoved(session)
			} else {
				Log.Warnf("session closed, reason: %s", msg.Reason)
			}
		case Msg:
			if session, ok := self.GetSession(ev.Session().ID()); ok {
				msg.SetSession(session)
				if dispatch, ok := msg.(*Dispatch); ok {
					dipHandler(dispatch)
				} else {
					evtHandler.OnRcvMessage(msg)
				}
			}
		}
	})

	if reconnector, ok := self.peer.(interface{
		SetReconnectDuration(time.Duration)
	}); ok {
		reconnector.SetReconnectDuration(5 * time.Second)
	}

	self.peer.Start()
}

func (self *CorePort) AddSession(id int64, session Session) {
	self.sLock.Lock()
	defer self.sLock.Unlock()
	self.sessions[id] = session
}

func (self *CorePort) RemoveSession(id int64) {
	self.sLock.Lock()
	defer self.sLock.Unlock()
	delete(self.sessions, id)
}

func (self *CorePort) GetSession(id int64) (session Session, ok bool) {
	self.sLock.RLock()
	defer self.sLock.RUnlock()
	session, ok = self.sessions[id]
	return
}

func (self *CorePort) Name() string {
	return self.name
}

func (self *CorePort) Stop() {
	self.pLock.Lock()
	defer self.pLock.Unlock()
	self.peer.Stop()
}

func (self *CorePort) SetService(service Service) {
	self.service = service
}

func (self *CorePort) Service() Service {
	return self.service
}

func (self *CorePort) Peer() cellnet.Peer {
	return self.peer
}