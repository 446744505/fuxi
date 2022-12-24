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
	Log.Debugf("role %v start enter map", self.RoleId)
	mapPvidId := MapMgr.RandomMap()
	if mapPvidId == 0 {
		Log.Errorf("role %v no map can use", self.RoleId)
		return
	}
	self.MapPvid = mapPvidId
	enter := &msg.GEnterMap{}
	enter.RoleId = self.RoleId
	enter.ClientSid = self.ClientSid
	enter.GsPvid = GS.Pvid
	enter.ProviderName = self.Provider.Port().HostPortString()
	if ok := self.SendToSelfMap(enter); !ok {
		Log.Errorf("role %v enter map failed", self.RoleId)
	}
}

func (self *NetRole) SendToSelfMap(msg core.Msg) bool {
	return GS.SendToProvidee(self.MapPvid, msg)
}
