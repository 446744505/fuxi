package provider

import (
	"fuxi/core"
	"fuxi/msg"
)

func OnDispatch(p *core.Dispatch) {
	if p.ToType() == core.MsgToProvidee {
		toPVID := int32(p.ToID())
		if sess, ok := Provider.providees[toPVID]; ok {
			sess.SendRaw(p.MsgId, p.MsgData)
		}
	} else if p.ToType() == core.MsgToClient {

	}
}

func (self *ProviderEventHandler) OnBindPvid(p core.Msg) {
	bind := p.(*msg.BindPvid)
	Provider.BindProvidee(bind.PVID, bind.Name, p.Session())
}

func (self *ProviderEventHandler) OnUnBindPvid(p core.Msg) {
	unbind := p.(*msg.UnBindPvid)
	Provider.UnBindProvidee(unbind.PVID)
}