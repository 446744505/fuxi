package linker

import (
	"fuxi/core"
	"fuxi/msg"
	"fuxi/switcher/util"
	"github.com/davyxu/cellnet"
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
	peer := p.Session().Port().Peer()
	ctx := peer.(cellnet.ContextSet)
	info, _ := ctx.GetContext(util.CtxTypeSessionInfo)
	linkerInfo := info.(*util.LinkerSessionInfo)
	linkerInfo.RoleId = center.RoleId

	lenter := &msg.LEnterGame{
		RoleId: center.RoleId,
		ClientSid: p.Session().ID(),
	}
	Log.Debugf("role %v enter", center.RoleId)
	util.SendToProvidee(center.PVID, lenter)
}