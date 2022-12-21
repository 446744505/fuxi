package internal

import (
	"fmt"
	"fuxi/core"
	"fuxi/msg"
	"fuxi/providee"
	"github.com/davyxu/golog"
	"strconv"
)

var GS *gs
var Log = golog.New("gs")

type gs struct {
	core.NetControlImpl

	Pvid int32
	roles map[int64]*NetRole
}

func NewGs() *gs {
	GS = &gs{
		roles: make(map[int64]*NetRole),
	}
	pvid, _ := strconv.Atoi(core.Args.Get("pvid"))
	GS.Pvid = int32(pvid)
	p := providee.NewProvidee(GS.Pvid, "gs")
	p.SetEventHandler(&gsEventHandler{})
	GS.AddService(p)

	core.ETCD.Put(fmt.Sprintf("gs/%v", pvid), fmt.Sprintf("%v", pvid))

	return GS
}

func (self *gs) OnRoleEnter(p *msg.LEnterGame) {
	role := &NetRole{}
	role.RoleId = p.RoleId
	role.ClientSid = p.ClientSid
	role.Provider = p.Session()
	role.EnterMap()
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
		role.Provider.Send(msg)
		return
	}
	Log.Errorf("not client role %d", roleId)
}

func (self *gs) SendToProvidee(pvid int32, msg core.Msg) {
	svr := self.GetService("gs")
	if p, ok := svr.(interface{
		SendToProvidee(int32, core.Msg)
	}); ok {
		p.SendToProvidee(pvid, msg)
	}
}