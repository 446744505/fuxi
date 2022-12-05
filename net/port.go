package net

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
)

type PortType = int

const (
	Acceptor PortType = iota
	Connector
)

type Porter interface {
	Start()
	Stop()
}

type Port struct {
	Conf *PortConf
	peerType string
	peer cellnet.Peer
}

type PortConf struct {
	*ServiceConf

	Typ PortType
	Host string
	Port int
	Ops *NetOptions
}

func NewPortConf(typ PortType) *PortConf {
	return &PortConf{Typ: typ}
}

func NewAccept(conf *PortConf) *Port {
	return &Port{
		Conf:     conf,
		peerType: "tcp.Acceptor",
	}
}

func NewConnect(conf *PortConf) *Port {
	return &Port{
		Conf:     conf,
		peerType: "tcp.Connector",
	}
}

func (self *Port) Start() {
	factory := self.Conf.ServiceConf.Factory
	handler := factory.NewEventHandler()
	self.peer = peer.NewGenericPeer(self.peerType, self.Conf.Name(), self.Conf.HostPortString(), nil)
	proc.BindProcessorHandler(self.peer, "tcp.ltv", func(ev cellnet.Event) {
		session := factory.NewSession(ev.Session())
		session.Port = self
		switch msg := ev.Message().(type) {
			case *cellnet.SessionAccepted, *cellnet.SessionConnected:
				handler.OnSessionAdd(session)
			case *cellnet.SessionClosed:
				handler.OnSessionRemoved(session)
			case *Msg:
				msg.Session = ev.Session()
				handler.OnRcvMessage(msg)
		}
	})
	self.peer.Start()
}

func (self *Port) Stop() {
	self.peer.Stop()
}

func (self *PortConf) BuildHost(host string) *PortConf {
	self.Host = host
	return self
}

func (self *PortConf) BuildPort(port int) *PortConf {
	self.Port = port
	return self
}

func (self *PortConf) HostPortString() string {
	return fmt.Sprintf("%s:%d", self.Host, self.Port)
}

func (self *PortConf) Name() string {
	return fmt.Sprintf("%s-%s", self.ServiceConf.Name, self.HostPortString())
}