package net

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"reflect"
)

type Ider interface {
	ID() int16
}

type Msg struct {
	Ider
	Session cellnet.Session
}

func RegisterMsg(msg Ider) {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf(msg).Elem(),
		ID:    int(msg.ID()),
	})
}