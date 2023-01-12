package core

type MsgHandler func(p Msg)

type EventHandler interface {
	Init()
	RegisterMsg(msg Msg, handler MsgHandler)
	OnSessionAdd(session Session)
	OnSessionRemoved(session Session)
	OnRcvMessage(msg Msg)
}

type CoreEventHandler struct {
	msgHandlers map[int32]MsgHandler
}

func (self *CoreEventHandler) Init() {
}

func (self *CoreEventHandler) OnSessionAdd(_ Session) {
}

func (self *CoreEventHandler) OnSessionRemoved(_ Session) {
}

func (self *CoreEventHandler) OnRcvMessage(msg Msg) {
	msg.Session().
		Port().Service().Pool().
		Process(func() {
			if handler, ok := self.msgHandlers[msg.ID()]; ok {
				handler(msg)
			}
		})
}

func (self *CoreEventHandler) OnRcvMessageSync(msg Msg) {
	if handler, ok := self.msgHandlers[msg.ID()]; ok {
		handler(msg)
	}
}

func (self *CoreEventHandler) RegisterMsg(msg Msg, handler MsgHandler) {
	if handler != nil {
		if self.msgHandlers == nil {
			self.msgHandlers = make(map[int32]MsgHandler)
		}
		self.msgHandlers[msg.ID()] = handler
	}
	RegisterMsg(msg)
}
