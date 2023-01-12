package msg

import (
	"fuxi/core"
	"github.com/davyxu/protoplus/proto"
)

type CEnterGame struct {
	core.CoreMsg `binary:"-"`
	RoleId       int64
	PVID         int32
}

func (self *CEnterGame) ID() int32 {
	return 5001
}

func (self *CEnterGame) Marshal(buffer *proto.Buffer) error {
	proto.MarshalInt64(buffer, 0, self.RoleId)
	proto.MarshalInt32(buffer, 1, self.PVID)
	return nil
}

func (self *CEnterGame) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalInt64(buffer, wt, &self.RoleId)
	case 1:
		return proto.UnmarshalInt32(buffer, wt, &self.PVID)
	}

	return proto.ErrUnknownField
}

func (self *CEnterGame) Size() (ret int) {
	ret += proto.SizeInt64(0, self.RoleId)
	ret += proto.SizeInt32(1, self.PVID)
	return
}

type SEnterGame struct {
	core.CoreMsg `binary:"-"`
	Name         string
}

func (self *SEnterGame) ID() int32 {
	return 5002
}

func (self *SEnterGame) Marshal(buffer *proto.Buffer) error {
	proto.MarshalString(buffer, 0, self.Name)
	return nil
}

func (self *SEnterGame) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalString(buffer, wt, &self.Name)
	}

	return proto.ErrUnknownField
}

func (self *SEnterGame) Size() (ret int) {
	ret += proto.SizeString(0, self.Name)
	return
}

type SEnterMap struct {
	core.CoreMsg `binary:"-"`
	Pvid         int32
}

func (self *SEnterMap) ID() int32 {
	return 5003
}

func (self *SEnterMap) Marshal(buffer *proto.Buffer) error {
	proto.MarshalInt32(buffer, 0, self.Pvid)
	return nil
}

func (self *SEnterMap) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalInt32(buffer, wt, &self.Pvid)
	}

	return proto.ErrUnknownField
}

func (self *SEnterMap) Size() (ret int) {
	ret += proto.SizeInt32(0, self.Pvid)
	return
}

type CGetInfo struct {
	core.CoreMsg `binary:"-"`
}

func (self *CGetInfo) ID() int32 {
	return 5004
}

func (self *CGetInfo) Marshal(buffer *proto.Buffer) error {
	return nil
}

func (self *CGetInfo) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	return nil
}

func (self *CGetInfo) Size() (ret int) {
	return 0
}
