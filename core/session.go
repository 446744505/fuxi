package core

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"net"
	"sync"
)

type SessionCreater func(raw cellnet.Session) Session

type Session interface {
	ID() int64
	Send(msg Msg)
	SendRaw(msgId int, msgData []byte)
	SetPort(port Port)
	Port() Port
	SetRaw(raw cellnet.Session)
	Raw() cellnet.Session
	SetContext(key string, val interface{})
	GetContext(key string) (interface{}, bool)
}

type CoreSession struct {
	raw cellnet.Session
	port Port

	lock sync.RWMutex
	ctx map[string]interface{}
}

func NewCoreSession(raw cellnet.Session) *CoreSession {
	return &CoreSession{
		raw: raw,
		ctx: make(map[string]interface{}),
	}
}

func (self *CoreSession) SetContext(key string, val interface{}) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.ctx[key] = val
}

func (self *CoreSession) GetContext(key string) (val interface{}, ok bool) {
	self.lock.Lock()
	defer self.lock.Unlock()
	val, ok = self.ctx[key]
	return
}


func (self *CoreSession) Send(msg Msg) {
	self.raw.Send(msg)
}

func (self *CoreSession) SendRaw(msgId int, msgData []byte) {
	self.raw.Send(&cellnet.RawPacket{
		MsgData: msgData,
		MsgID: msgId,
	})
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

func (self *CoreSession) Raw() cellnet.Session {
	return self.raw
}

func (self *CoreSession) SetRaw(raw cellnet.Session) {
	self.raw = raw
}

func (self *CoreSession) String() string {
	conn := self.raw.Raw().(*net.TCPConn)
	return fmt.Sprintf("%s->%s", conn.LocalAddr().String(), conn.RemoteAddr().String())
}