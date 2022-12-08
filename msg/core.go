package msg

import (
	"fuxi/core"
	"github.com/davyxu/protoplus/proto"
)

type Dispatch struct {
	core.CoreMsg `binary:"-"`
	ClientSid int64
	MsgData []byte
}

func (self Dispatch) ID() int32 {
	return 1
}

func (self *Dispatch) Marshal(buffer *proto.Buffer) error {
	proto.MarshalInt64(buffer, 0, self.ClientSid)
	proto.MarshalBytes(buffer, 1, self.MsgData)
	return nil
}

func (self *Dispatch) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalInt64(buffer, wt, &self.ClientSid)
	case 1:
		return proto.UnmarshalBytes(buffer, wt, &self.MsgData)
	}

	return proto.ErrUnknownField
}

func (self *Dispatch) Size() (ret int) {
	ret += proto.SizeInt64(0, self.ClientSid)
	ret += proto.SizeBytes(1, self.MsgData)
	return
}

type PDispatch struct {
	core.CoreMsg `binary:"-"`
	PVID int32
	MsgData []byte
}

func (self PDispatch) ID() int32 {
	return 2
}

func (self *PDispatch) Marshal(buffer *proto.Buffer) error {
	proto.MarshalInt32(buffer, 0, self.PVID)
	proto.MarshalBytes(buffer, 1, self.MsgData)
	return nil
}

func (self *PDispatch) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalInt32(buffer, wt, &self.PVID)
	case 1:
		return proto.UnmarshalBytes(buffer, wt, &self.MsgData)
	}

	return proto.ErrUnknownField
}

func (self *PDispatch) Size() (ret int) {
	ret += proto.SizeInt32(0, self.PVID)
	ret += proto.SizeBytes(1, self.MsgData)
	return
}

type BindPvid struct {
	core.CoreMsg `binary:"-"`
	PVID int32
	Name string
}

func (self BindPvid) ID() int32 {
	return 3
}

func (self *BindPvid) Marshal(buffer *proto.Buffer) error {
	proto.MarshalInt32(buffer, 0, self.PVID)
	proto.MarshalString(buffer, 1, self.Name)
	return nil
}

func (self *BindPvid) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalInt32(buffer, wt, &self.PVID)
	case 1:
		return proto.UnmarshalString(buffer, wt, &self.Name)
	}

	return proto.ErrUnknownField
}

func (self *BindPvid) Size() (ret int) {
	ret += proto.SizeInt32(0, self.PVID)
	ret += proto.SizeString(1, self.Name)
	return
}

type UnBindPvid struct {
	core.CoreMsg `binary:"-"`
	PVID int32
}

func (self UnBindPvid) ID() int32 {
	return 4
}

func (self *UnBindPvid) Marshal(buffer *proto.Buffer) error {
	proto.MarshalInt32(buffer, 0, self.PVID)
	return nil
}

func (self *UnBindPvid) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalInt32(buffer, wt, &self.PVID)
	}

	return proto.ErrUnknownField
}

func (self *UnBindPvid) Size() (ret int) {
	ret += proto.SizeInt32(0, self.PVID)
	return
}