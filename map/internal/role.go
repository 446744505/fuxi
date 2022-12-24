package internal

import (
	"fuxi/core"
	"fuxi/msg"
)

type NetRole struct {
	*msg.GEnterMap
}

func (self *NetRole) OnEnterMap(enter *msg.GEnterMap) {
	self.GEnterMap = enter
	Log.Infof("role %v enter map", enter.RoleId)
	if ok := self.Send(&msg.SEnterMap{Pvid: Map.Pvid}); !ok {
		Log.Errorf("role %v map notify client failed", enter.RoleId)
	}
}

func (self *NetRole) Send(msg core.Msg) bool {
	return Map.SendToClient(self.ProviderName, self.ClientSid, msg)
}

func (self *NetRole) SendToGs(msg core.Msg) bool {
	return Map.SendToProvidee(self.GsPvid, msg)
}

