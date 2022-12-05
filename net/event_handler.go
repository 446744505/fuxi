package net

type EventHandler interface {
	OnSessionAdd(session *Session)
	OnSessionRemoved(session *Session)
	OnRcvMessage(msg *Msg)
}

type SampleEventHandler struct {

}

func (self SampleEventHandler) OnSessionAdd(session *Session) {

}

func (self SampleEventHandler) OnSessionRemoved(session *Session) {

}

func (self SampleEventHandler) OnRcvMessage(msg *Msg) {

}

