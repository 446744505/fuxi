package core

import (
	"github.com/davyxu/cellnet"
	"io"
	"net"
)

type CoreMessageTransmitter struct {
}

type socketOpt interface {
	MaxPacketSize() int
	ApplySocketReadTimeout(conn net.Conn, callback func())
	ApplySocketWriteTimeout(conn net.Conn, callback func())
}

func (self CoreMessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {
	reader, ok := ses.Raw().(io.Reader)

	// 转换错误，或者连接已经关闭时退出
	if !ok || reader == nil {
		return nil, nil
	}

	opt := ses.Peer().(socketOpt)

	if conn, ok := reader.(net.Conn); ok {

		// 有读超时时，设置超时
		opt.ApplySocketReadTimeout(conn, func() {

			msg, err = RecvLTVPacket(reader, opt.MaxPacketSize())
			if err != nil {
				Log.Errorf("receive msg err: %v", err)
			}
		})
	}

	return
}

func (self CoreMessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) (err error) {
	writer, ok := ses.Raw().(io.Writer)

	// 转换错误，或者连接已经关闭时退出
	if !ok || writer == nil {
		return nil
	}

	opt := ses.Peer().(socketOpt)

	// 有写超时时，设置超时
	opt.ApplySocketWriteTimeout(writer.(net.Conn), func() {

		err = SendLTVPacket(writer, ses.(cellnet.ContextSet), msg)
		if err != nil {
			Log.Errorf("send msg err: %v", err)
		}
	})

	return
}
