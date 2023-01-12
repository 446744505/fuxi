package core

import "fmt"

type ServiceConf interface {
	SetName(name string)
	Name() string
	SetPoolCapacity(capacity int)
	PoolCapacity() int
}

type CoreServiceConf struct {
	name         string
	poolCapacity int
}

func (self *CoreServiceConf) Name() string {
	return self.name
}

func (self *CoreServiceConf) SetName(name string) {
	self.name = name
}

func (self *CoreServiceConf) PoolCapacity() int {
	return self.poolCapacity
}

func (self *CoreServiceConf) SetPoolCapacity(capacity int) {
	self.poolCapacity = capacity
}

type PortType = string

const (
	TypeAcceptor  PortType = "tcp.Acceptor"
	TypeConnector PortType = "tcp.Connector"
)

type CorePortConf struct {
	typ  PortType
	name string
	host string
	port int
}

func (self *CorePortConf) HostPortString() string {
	return fmt.Sprintf("%s:%d", self.host, self.port)
}
