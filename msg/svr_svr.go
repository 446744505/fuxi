package msg

import (
	"fuxi/core"
	"github.com/davyxu/protoplus/proto"
)

type MapNtf struct {
	core.CoreMsg `binary:"-"`
	PVID int32
}

func (self *MapNtf) ID() int32 {
	return 1001
}

func (self *MapNtf) Marshal(buffer *proto.Buffer) error {
	proto.MarshalInt32(buffer, 0, self.PVID)
	return nil
}

func (self *MapNtf) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalInt32(buffer, wt, &self.PVID)
	}

	return proto.ErrUnknownField
}

func (self *MapNtf) Size() (ret int) {
	ret += proto.SizeInt32(0, self.PVID)
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
	ret += proto.SizeInt64(0, self.ClientSid)
	return
}