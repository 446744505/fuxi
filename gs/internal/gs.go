package internal

import (
	"fuxi/core"
	"fuxi/msg"
	"fuxi/providee"
	"github.com/davyxu/golog"
)

var GS *gs
var Log = golog.New("gs")

type gs struct {
	core.NetControlImpl

	roles map[int64]*NetRole
}

func NewGs() *gs {
	GS = &gs{
		roles: make(map[int64]*NetRole),
	}
	p := providee.NewProvidee(1, "gs")
	p.SetEventHandler(&gsEventHandler{})
	GS.AddService(p)
	return GS
}

func (self *gs) OnRoleEnter(p *msg.LEnterGame) {
	role := &NetRole{}
	role.ClientSid = p.ClientSid
	role.Session = p.Session()
	self.roles[p.RoleId] = role
	ack := &msg.SEnterGame{}
	ack.Name = "玩家1"
	self.SendToClient(p.RoleId, ack)
	Log.Infof("role %d enter game", p.RoleId)
}

func (self *gs) SendToClient(roleId int64, msg core.Msg) {
	if role, ok := self.roles[roleId]; ok {
		msg.SetToType(core.MsgToClient)
		msg.SetToID(role.ClientSid)
		role.Session.Send(msg)
		return
	}
	Log.Errorf("not client role %d", roleId)
}