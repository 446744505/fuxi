package linker

import (
	"fuxi/core"
	"fuxi/msg"
	"fuxi/switcher/util"
)

func OnDispatch(p *core.Dispatch) {
	if p.ToType() == core.MsgToServer {
		util.ClientToProvidee(p)
	} else {
		Log.Errorln("err dispatch msg, toType: %d msgId: %d", p.ToType(), p.MsgId)
	}
}

func (self *linkerEventHandler) OnCEnterGame(p core.Msg) {
	center := p.(*msg.CEnterGame)
	info, _ := p.Session().GetContext(util.CtxTypeSessionInfo)
	linkerInfo := info.(*util.LinkerSessionInfo)
	linkerInfo.RoleId = center.RoleId

	lenter := &msg.LEnterGame{
		RoleId: center.RoleId,
		ClientSid: p.Session().ID(),
	}
	Log.Debugf("role %v enter", center.RoleId)
	util.SendToProvidee(center.PVID, lenter)
}