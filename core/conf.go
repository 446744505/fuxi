package core

import "fmt"

type ServiceConf interface {
	SetName(name string)
	Name() string
}

type CoreServiceConf struct {
	name string
}

func (self *CoreServiceConf) Name() string {
	return self.name
}

func (self *CoreServiceConf) SetName(name string) {
	self.name = name
}

type PortType = string

const (
	TypeAcceptor PortType = "tcp.Acceptor"
	TypeConnector PortType = "tcp.Connector"
)

type PortConf interface {
	GetPortType() PortType
	SetPortType(typ PortType)
	Host() string
	SetHost(host string)
	Port() int
	SetPort(port int)
}

type CorePortConf struct {
	typ PortType
	host string
	port int
}

func (self *CorePortConf) GetPortType() PortType {
	return self.typ
}

func (self *CorePortConf) SetPortType(typ PortType) {
	self.typ = typ
}

func (self *CorePortConf) Host() string {
	return self.host
}

func (self *CorePortConf) SetHost(host string) {
	self.host = host
}

func (self *CorePortConf) Port() int {
	return self.port
}

func (self *CorePortConf) SetPort(port int) {
	self.port = port
}

func (self *CorePortConf) HostPortString() string {
	return fmt.Sprintf("%s:%d", self.host, self.port)
}

func NewAcceptor(host string, port int) Port {
	t := &CorePort{}
	t.typ = TypeAcceptor
	t.host = host
	t.port = port
	return t
}

func NewConnector(host string, port int) Port {
	t := &CorePort{}
	t.typ = TypeConnector
	t.host = host
	t.port = port
	return t
}