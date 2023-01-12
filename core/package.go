package core

import (
	"encoding/binary"
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/util"
	"io"
)

var (
	ErrUnhandledMsg = errors.New("unhandled msg")
	ErrMaxPacket    = errors.New("packet over size")
	ErrMinPacket    = errors.New("packet short size")
	ErrShortMsgInfo = errors.New("short msginfo")
)

const (
	bodySize      = 2 // 包体大小字段
	msgToTypeSize = 2
	msgToIDSize   = 8
	msgIDSize     = 2 // 消息ID字段
	msgInfoSize   = msgToTypeSize + msgToIDSize + msgIDSize
)

// 接收Length-Type-Value格式的封包流程
func RecvLTVPacket(reader io.Reader, maxPacketSize int) (msg interface{}, err error) {

	// Size为uint16，占2字节
	var sizeBuffer = make([]byte, bodySize)

	// 持续读取Size直到读到为止
	_, err = io.ReadFull(reader, sizeBuffer)

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(sizeBuffer) < bodySize {
		return nil, ErrMinPacket
	}

	// 用小端格式读取Size
	size := binary.LittleEndian.Uint16(sizeBuffer)

	if maxPacketSize > 0 && size >= uint16(maxPacketSize) {
		return nil, ErrMaxPacket
	}

	// 分配包体大小
	body := make([]byte, size)

	// 读取包体数据
	_, err = io.ReadFull(reader, body)

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(body) < msgInfoSize {
		return nil, ErrShortMsgInfo
	}

	msgToType := binary.LittleEndian.Uint16(body)

	msgToId := binary.LittleEndian.Uint64(body[msgToTypeSize:])

	msgid := binary.LittleEndian.Uint16(body[msgToTypeSize+msgToIDSize:])

	msgData := body[msgInfoSize:]

	mid := int(msgid)
	meta := cellnet.MessageMetaByID(mid)
	if meta == nil {
		if msgToType == 0 {
			return nil, ErrUnhandledMsg
		}
		return NewDispatch(MsgToType(msgToType), int64(msgToId), mid, msgData), nil
	}

	// 将字节数组和消息ID用户解出消息
	msg, _, err = codec.DecodeMessage(mid, msgData)
	if err != nil {
		return nil, err
	}

	return
}

// 发送Length-Type-Value格式的封包流程
func SendLTVPacket(writer io.Writer, ctx cellnet.ContextSet, data interface{}) error {

	var (
		msgData   []byte
		msgID     int
		meta      *cellnet.MessageMeta
		msgToType int16
		msgToId   int64
	)

	switch m := data.(type) {
	case *cellnet.RawPacket: // 发裸包
		msgData = m.MsgData
		msgID = m.MsgID
	default: // 发普通编码包
		var err error

		// 将用户数据转换为字节数组和消息ID
		msgData, meta, err = codec.EncodeMessage(data, ctx)

		if err != nil {
			return err
		}

		msgID = meta.ID
	}

	if msg, ok := data.(Msg); ok {
		msgToType = int16(msg.ToType())
		msgToId = msg.FTId()
	}

	pkt := make([]byte, bodySize+msgInfoSize+len(msgData))

	// Length
	binary.LittleEndian.PutUint16(pkt, uint16(msgInfoSize+len(msgData)))

	//to type
	binary.LittleEndian.PutUint16(pkt[bodySize:], uint16(msgToType))

	//to id
	binary.LittleEndian.PutUint64(pkt[bodySize+msgToTypeSize:], uint64(msgToId))

	// Type
	binary.LittleEndian.PutUint16(pkt[bodySize+msgToTypeSize+msgToIDSize:], uint16(msgID))

	// Value
	copy(pkt[bodySize+msgInfoSize:], msgData)

	// 将数据写入Socket
	err := util.WriteFull(writer, pkt)

	// Codec中使用内存池时的释放位置
	if meta != nil {
		codec.FreeCodecResource(meta.Codec, msgData, ctx)
	}

	return err
}
