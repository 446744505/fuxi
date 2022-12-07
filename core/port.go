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
	handler := self.Service().EventHandler()
	proc.BindProcessorHandler(self.peer, "tcp.ltv", func(ev cellnet.Event) {
		ctx := self.peer.(cellnet.ContextSet)
		switch msg := ev.Message().(type) {
		case *cellnet.SessionAccepted, *cellnet.SessionConnected:
			session := NewSession(ev.Session())
			session.SetPort(self)
			handler.OnSessionAdd(session)
			ctx.SetContext(CtxTypeSession, session)
		case *cellnet.SessionClosed:
			if val, ok := ctx.GetContext(CtxTypeSession); ok {
				session := val.(Session)
				handler.OnSessionRemoved(session)
			}
		case Msg:
			if val, ok := ctx.GetContext(CtxTypeSession); ok {
				session := val.(Session)
				msg.SetSession(session)
				handler.OnRcvMessage(msg)
			}
		}
	})
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