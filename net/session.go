package net

import "github.com/davyxu/cellnet"

type Session struct {
	*Port
	Raw cellnet.Session
}

func NewSession(raw cellnet.Session) *Session {
	return &Session{Raw: raw}
}

func (self *Session) Send(msg *Msg) {
	self.Raw.Send(msg)
}
