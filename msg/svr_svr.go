package msg

import (
	"fuxi/core"
	"github.com/davyxu/protoplus/proto"
)

type GEnterMap struct {
	core.CoreMsg `binary:"-"`
	RoleId int64
	ClientSid int64
	GsPvid int32
	ProviderName string
}

func (self *GEnterMap) ID() int32 {
	return 1001
}

func (self *GEnterMap) Marshal(buffer *proto.Buffer) error {
	proto.MarshalInt64(buffer, 0, self.RoleId)
	proto.MarshalInt64(buffer, 1, self.ClientSid)
	proto.MarshalInt32(buffer, 2, self.GsPvid)
	proto.MarshalString(buffer, 3, self.ProviderName)
	return nil
}

func (self *GEnterMap) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalInt64(buffer, wt, &self.RoleId)
	case 1:
		return proto.UnmarshalInt64(buffer, wt, &self.ClientSid)
	case 2:
		return proto.UnmarshalInt32(buffer, wt, &self.GsPvid)
	case 3:
		return proto.UnmarshalString(buffer, wt, &self.ProviderName)
	}

	return proto.ErrUnknownField
}

func (self *GEnterMap) Size() (ret int) {
	ret += proto.SizeInt64(0, self.RoleId)
	ret += proto.SizeInt64(1, self.ClientSid)
	ret += proto.SizeInt32(2, self.GsPvid)
	ret += proto.SizeString(3, self.ProviderName)
	return
}

type LEnterGame struct {
	core.CoreMsg `binary:"-"`
	RoleId int64
	ClientSid int64
}

func (self *LEnterGame) ID() int32 {
	return 1002
}

func (self *LEnterGame) Marshal(buffer *proto.Buffer) error {
	proto.MarshalInt64(buffer, 0, self.RoleId)
	proto.MarshalInt64(buffer, 1, self.ClientSid)
	return nil
}

func (self *LEnterGame) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalInt64(buffer, wt, &self.RoleId)
	case 1:
		return proto.UnmarshalInt64(buffer, wt, &self.ClientSid)
	}

	return proto.ErrUnknownField
}

func (self *LEnterGame) Size() (ret int) {
	ret += proto.SizeInt64(0, self.RoleId)
	ret += proto.SizeInt64(1, self.ClientSid)
	return
}

type MessageBox struct {
	core.CoreMsg `binary:"-"`
	UniqId       int64
	MsgId        int32
	MsgData      []byte
}

func (self *MessageBox) ID() int32 {
	return 1003
}

func (self *MessageBox) Marshal(buffer *proto.Buffer) error {
	proto.MarshalInt64(buffer, 0, self.UniqId)
	proto.MarshalInt32(buffer, 1, self.MsgId)
	proto.MarshalBytes(buffer, 2, self.MsgData)
	return nil
}

func (self *MessageBox) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalInt64(buffer, wt, &self.UniqId)
	case 1:
		return proto.UnmarshalInt32(buffer, wt, &self.MsgId)
	case 2:
		return proto.UnmarshalBytes(buffer, wt, &self.MsgData)
	}

	return proto.ErrUnknownField
}

func (self *MessageBox) Size() (ret int) {
	ret += proto.SizeInt64(0, self.UniqId)
	ret += proto.SizeInt32(1, self.MsgId)
	ret += proto.SizeBytes(2, self.MsgData)
	return
}

type ClientBroken struct {
	core.CoreMsg `binary:"-"`
	ClientSid       int64
}

func (self *ClientBroken) ID() int32 {
	return 1004
}

func (self *ClientBroken) Marshal(buffer *proto.Buffer) error {
	proto.MarshalInt64(buffer, 0, self.ClientSid)
	return nil
}

func (self *ClientBroken) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalInt64(buffer, wt, &self.ClientSid)
	}

	return proto.ErrUnknownField
}

func (self *ClientBroken) Size() (ret int) {
	ret += proto.SizeInt64(0, self.ClientSid)
	return
}

type GExitMap struct {
	core.CoreMsg `binary:"-"`
	RoleId int64
}

func (self *GExitMap) ID() int32 {
	return 1005
}

func (self *GExitMap) Marshal(buffer *proto.Buffer) error {
	proto.MarshalInt64(buffer, 0, self.RoleId)
	return nil
}

func (self *GExitMap) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalInt64(buffer, wt, &self.RoleId)
	}

	return proto.ErrUnknownField
}

func (self *GExitMap) Size() (ret int) {
	ret += proto.SizeInt64(0, self.RoleId)
	return
}