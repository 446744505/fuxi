package core

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
	"sync"
	"sync/atomic"
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
	peer atomic.Value

	sessions sync.Map
}

func (self *CorePort) Start() {
	peer := peer.NewGenericPeer(self.typ, self.name, fmt.Sprintf("0.0.0.0:%d", self.port), nil)
	self.peer.Store(peer)
	evtHandler := self.Service().EventHandler()
	creater := self.Service().SessionCreater()
	dipHandler := self.Service().DispatchHandler()
	proc.BindProcessorHandler(peer, "fxtcp.ltv", func(ev cellnet.Event) {
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

	if reconnector, ok := peer.(interface{
		SetReconnectDuration(time.Duration)
	}); ok {
		reconnector.SetReconnectDuration(5 * time.Second)
	}

	peer.Start()
}

func (self *CorePort) AddSession(id int64, session Session) {
	self.sessions.Store(id, session)
}

func (self *CorePort) RemoveSession(id int64) {
	self.sessions.Delete(id)
}

func (self *CorePort) GetSession(id int64) (Session, bool) {
	if val, ok := self.sessions.Load(id); ok {
		return val.(Session), true
	}
	return nil, false
}

func (self *CorePort) Name() string {
	return self.name
}

func (self *CorePort) Stop() {
	self.Peer().Stop()
}

func (self *CorePort) SetService(service Service) {
	self.service = service
}

func (self *CorePort) Service() Service {
	return self.service
}

func (self *CorePort) Peer() cellnet.Peer {
	return self.peer.Load().(cellnet.Peer)
}