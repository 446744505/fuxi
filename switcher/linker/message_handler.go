package linker

import (
	"fuxi/core"
	"fuxi/msg"
	"fuxi/switcher/linker/util"
)

func OnDispatch(p *core.Dispatch) {
	if p.ToType() == core.MsgToServer {
		util.DispatchToProvidee(p)
	} else {
		Log.Errorln("err dispatch msg, toType: %d msgId: %d", p.ToType(), p.MsgId)
	}
}

func (self *linkerEventHandler) OnCEnterGame(p core.Msg) {
	center := p.(*msg.CEnterGame)
	lenter := &msg.LEnterGame{
		RoleId: center.RoleId,
		ClientSid: p.Session().ID(),
	}
	util.SendToProvidee(center.PVID, lenter)
}