package core

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
)

type CtxType int

const (
	CtxTypeSession = 1
)

type Port interface {
	Controler
	SetService(service Service)
	Service() Service
}

type CorePort struct {
	CorePortConf
	service Service
	peer cellnet.Peer
}

func (self *CorePort) Start() {
	self.peer = peer.NewGenericPeer(self.typ, "", self.HostPortString(), nil)
	evtHandler := self.Service().EventHandler()
	creater := self.Service().SessionCreater()
	dipHandler := self.Service().DispatchHandler()
	proc.BindProcessorHandler(self.peer, "fxtcp.ltv", func(ev cellnet.Event) {
		ctx := self.peer.(cellnet.ContextSet)
		switch msg := ev.Message().(type) {
		case *cellnet.SessionAccepted, *cellnet.SessionConnected:
			session := creater(ev.Session())
			session.SetPort(self)
			evtHandler.OnSessionAdd(session)
			ctx.SetContext(CtxTypeSession, session)
		case *cellnet.SessionClosed:
			if val, ok := ctx.GetContext(CtxTypeSession); ok {
				session := val.(Session)
				Log.Infof("session %s closed, reason: %s", session, msg.Reason)
				evtHandler.OnSessionRemoved(session)
			} else {
				Log.Warnf("session closed, reason: %s", msg.Reason)
			}
		case Msg:
			if val, ok := ctx.GetContext(CtxTypeSession); ok {
				session := val.(Session)
				msg.SetSession(session)
				if dispatch, ok := msg.(*Dispatch); ok {
					dipHandler(dispatch)
				} else {
					evtHandler.OnRcvMessage(msg)
				}
			}
		}
	})

	/* 要重连就放开这里
	if reconnector, ok := self.peer.(interface{
		SetReconnectDuration(time.Duration)
	}); ok {
		reconnector.SetReconnectDuration(5 * time.Second)
	}
	 */

	self.peer.Start()
}

func (self *CorePort) Stop() {
	self.peer.Stop()
}

func (self *CorePort) SetService(service Service) {
	self.service = service
}

func (self *CorePort) Service() Service {
	return self.service
}