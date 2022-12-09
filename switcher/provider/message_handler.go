package provider

import (
	"fuxi/core"
	"fuxi/msg"
	"fuxi/switcher/linker/util"
)

func OnDispatch(p *core.Dispatch) {
	if p.ToType() == core.MsgToProvidee {
		util.DispatchToProvidee(p)
	} else if p.ToType() == core.MsgToClient {
		util.DispatchToClient(p)
	} else {
		Log.Errorln("err dispatch msg, toType: %d msgId: %d", p.ToType(), p.MsgId)
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