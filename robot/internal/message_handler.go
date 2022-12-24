package internal

import (
	"fuxi/core"
	"fuxi/msg"
	"github.com/davyxu/cellnet"
)

func GetRole(p core.Msg) *Role {
	ctx := p.Session().Port().Peer().(cellnet.ContextSet)
	if role, ok := ctx.GetContext(CtxRole); ok {
		return role.(*Role)
	}
	Log.Errorf("no role at session %v", p.Session())
	return nil
}

func (self *robotEventHandler) OnSEnterGame(p core.Msg) {
	enter := p.(*msg.SEnterGame)
	role := GetRole(p)
	role.EnterGame(enter)
}

func (self *robotEventHandler) OnSEnterMap(p core.Msg) {
	enter := p.(*msg.SEnterMap)
	role := GetRole(p)
	role.EnterMap(enter)
}