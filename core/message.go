package core

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/protoplus/proto"
	"reflect"
)

type MsgToType int16

const (
	MsgToProvidee MsgToType = 1
	MsgToClient MsgToType = 2
	MsgToServer MsgToType = 3
)

type MsgHead interface {
	SetToType(typ MsgToType)
	SetToID(id int64)
	ToType() MsgToType
	ToID() int64
}

type CoreMsgHead struct {
	toType MsgToType
	toID int64
}

type Msg interface {
	MsgHead
	ID() int32
	SetSession(session Session)
	Session() Session
	NewHead() MsgHead
}

type CoreMsg struct {
	session Session `binary:"-"`
	CoreMsgHead
}

func (self *CoreMsg) Session() Session {
	return self.session
}

func (self *CoreMsg) SetSession(session Session) {
	self.session = session
}

func (self *CoreMsg) NewHead() MsgHead {
	self.CoreMsgHead = CoreMsgHead{}
	return &self.CoreMsgHead
}

func (self *CoreMsgHead) SetToType(typ MsgToType) {
	self.toType = typ
}

func (self *CoreMsgHead) SetToID(id int64) {
	self.toID = id
}

func (self *CoreMsgHead) ToType() MsgToType {
	return self.toType
}

func (self *CoreMsgHead) ToID() int64 {
	return self.toID
}

func RegisterMsg(msg Msg) {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("protoplus"),
		Type:  reflect.TypeOf(msg).Elem(),
		ID:    int(msg.ID()),
	})
}

func NewDispatch(msgToType MsgToType, msgToId int64, msgId int, msgData []byte) interface{} {
	d := &Dispatch{}
	d.toType = msgToType
	d.toID = msgToId
	d.MsgData = msgData
	d.MsgId = msgId
	return d
}

type Dispatch struct {
	CoreMsg `binary:"-"`
	MsgId int
	MsgData []byte
}

func (self Dispatch) ID() int32 {
	return 1
}

func (self *Dispatch) Marshal(buffer *proto.Buffer) error {
	proto.MarshalBytes(buffer, 1, self.MsgData)
	return nil
}

func (self *Dispatch) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalBytes(buffer, wt, &self.MsgData)
	}

	return proto.ErrUnknownField
}

func (self *Dispatch) Size() (ret int) {
	ret += proto.SizeBytes(0, self.MsgData)
	return
}