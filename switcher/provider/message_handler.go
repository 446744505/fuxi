package provider

import (
	"fuxi/core"
	"fuxi/msg"
	"fuxi/switcher/util"
)

func OnDispatch(p *core.Dispatch) {
	if p.ToType() == core.MsgToProvidee {
		util.ProvideeToProvidee(p)
	} else if p.ToType() == core.MsgToClient {
		util.DispatchToClient(p)
	} else {
		Log.Errorf("err dispatch msg, toType: %d msgId: %d", p.ToType(), p.MsgId)
	}
}

func (self *providerEventHandler) OnBindPvid(p core.Msg) {
	bind := p.(*msg.BindPvid)
	Provider.BindProvidee(bind.PVID, bind.Name, p.Session())
}

func (self *providerEventHandler) OnUnBindPvid(p core.Msg) {
	unbind := p.(*msg.UnBindPvid)
	Provider.UnBindProvidee(unbind.PVID)
}
