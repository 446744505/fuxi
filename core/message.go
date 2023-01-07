package core

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/protoplus/proto"
	"reflect"
)

type MsgToType int16

const (
	MsgToProvidee MsgToType = 1 //服务器之间
	MsgToClient MsgToType = 2 //服务器到客户端
	MsgToServer MsgToType = 3 //客户端到服务器
)

type MsgHead interface {
	SetToType(typ MsgToType)
	SetFTId(id int64)
	ToType() MsgToType
	FTId() int64
}

type CoreMsgHead struct {
	toType MsgToType
	ftId   int64
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

func (self *CoreMsgHead) SetFTId(id int64) {
	self.ftId = id
}

func (self *CoreMsgHead) ToType() MsgToType {
	return self.toType
}

func (self *CoreMsgHead) FTId() int64 {
	return self.ftId
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
	d.ftId = msgToId
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
	return -1
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