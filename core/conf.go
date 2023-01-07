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

type CorePortConf struct {
	typ PortType
	name string
	host string
	port int
}

func (self *CorePortConf) HostPortString() string {
	return fmt.Sprintf("%s:%d", self.host, self.port)
}

func NewAcceptor(name, host string, port int) Port {
	t := &CorePort{}
	t.typ = TypeAcceptor
	t.name = name
	t.host = host
	t.port = port
	t.sessions = make(map[int64]Session)
	return t
}

func NewConnector(name, host string, port int) Port {
	t := &CorePort{}
	t.typ = TypeConnector
	t.name = name
	t.host = host
	t.port = port
	t.sessions = make(map[int64]Session)
	return t
}