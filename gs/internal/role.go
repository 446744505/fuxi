package internal

import (
	"fuxi/core"
	"fuxi/msg"
)

type NetRole struct {
	RoleId    int64
	ClientSid int64
	Provider  core.Session

	MapPvid int32
}

func (self *NetRole) EnterMap() {
	self.MapPvid = MapMgr.RandomMap()
	enter := &msg.GEnterMap{}
	enter.RoleId = self.RoleId
	enter.ClientSid = self.ClientSid
	enter.GsPvid = GS.Pvid
	self.SendToSelfMap(enter)
}

func (self *NetRole) SendToSelfMap(msg core.Msg) {
	GS.SendToProvidee(self.MapPvid, msg)
}
