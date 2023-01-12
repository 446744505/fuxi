package providee

import (
	"fuxi/core"
	"fuxi/msg"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

func (self *ProvideeEventHandler) OnMessageBox(p core.Msg) {
	box := p.(*msg.MessageBox)
	mid := int(box.MsgId)
	meta := cellnet.MessageMetaByID(mid)
	if meta == nil {
		core.Log.Errorf("MessageBox no msg %v", mid)
		return
	}
	m, _, err := codec.DecodeMessage(mid, box.MsgData)
	if err != nil {
		core.Log.Errorf("MessageBox decode msg err %v", err)
		return
	}
	msg := m.(core.Msg)
	msg.SetSession(p.Session())
	msg.SetFTId(box.UniqId)
	self.OnRcvMessageSync(msg)
}

func (self *ProvideeEventHandler) OnClientBroken(p core.Msg) {
	broken := p.(*msg.ClientBroken)
	onBroken := Providee.OnClientBroken
	if onBroken != nil {
		onBroken(broken.ClientSid)
	}
}