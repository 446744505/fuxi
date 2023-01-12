package internal

import (
	"fuxi/core"
	"fuxi/msg"
)

type Role struct {
	*msg.GEnterMap
}

func (self *Role) OnEnterMap(enter *msg.GEnterMap) {
	self.GEnterMap = enter
	Log.Infof("role %v enter map", enter.RoleId)
	if ok := self.Send(&msg.SEnterMap{Pvid: Map.pvid}); !ok {
		Log.Errorf("role %v map notify client failed", enter.RoleId)
	}
}

func (self *Role) OnExitMap() {
	Log.Infof("role %v exit map", self.RoleId)
}

func (self *Role) Send(msg core.Msg) bool {
	return Map.SendToClient(self.ProviderName, self.ClientSid, msg)
}

func (self *Role) SendToGs(msg core.Msg) bool {
	return Map.SendToProvidee(self.GsPvid, msg)
}
