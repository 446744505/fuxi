package core

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"reflect"
)

type Msg interface {
	ID() int32
	SetSession(session Session)
	Session() Session
}

type CoreMsg struct {
	session Session `binary:"-"`
}

func (self *CoreMsg) Session() Session {
	return self.session
}

func (self *CoreMsg) SetSession(session Session) {
	self.session = session
}

func RegisterMsg(msg Msg) {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("protoplus"),
		Type:  reflect.TypeOf(msg).Elem(),
		ID:    int(msg.ID()),
	})
}