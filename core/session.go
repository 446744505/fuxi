package core

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"net"
)

type Session interface {
	ID() int64
	Send(msg Msg)
	SetPort(port Port)
	Port() Port
}

type CoreSession struct {
	raw cellnet.Session
	port Port
}

func NewSession(raw cellnet.Session) Session {
	return &CoreSession{raw: raw}
}

func (self *CoreSession) Send(msg Msg) {
	self.raw.Send(msg)
}

func (self *CoreSession) ID() int64 {
	return self.raw.ID()
}

func (self *CoreSession) SetPort(port Port) {
	self.port = port
}

func (self *CoreSession) Port() Port {
	return self.port
}

func (self *CoreSession) String() string {
	conn := self.raw.Raw().(*net.TCPConn)
	return fmt.Sprintf("%s->%s", conn.LocalAddr().String(), conn.RemoteAddr().String())
}