package internal

import (
	"fmt"
	"fuxi/core"
	"fuxi/msg"
)

type Role struct {
	RoleId    int64
	ClientSid int64
	Provider  core.Session

	MapPvid int32
}

func (self *Role) EnterGame() {
	self.enterMap()

	ack := &msg.SEnterGame{}
	ack.Name = fmt.Sprintf("玩家%v", self.RoleId)
	self.Send(ack)
	Log.Infof("role %d enter game", self.RoleId)
}

func (self *Role) enterMap() {
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
	enter.GsPvid = GS.pvid
	enter.ProviderName = self.Provider.Port().HostPortString()
	if ok := self.SendToSelfMap(enter); !ok {
		Log.Errorf("role %v enter map failed", self.RoleId)
	}
}

func (self *Role) ExitGame() {
	defer core.PrintPanicStack()
	self.exitMap()
	Log.Infof("role %d exit game", self.RoleId)
}

func (self *Role) exitMap() {
	self.SendToSelfMap(&msg.GExitMap{RoleId: self.RoleId})
}

func (self *Role) Send(msg core.Msg) {
	msg.SetToType(core.MsgToClient)
	msg.SetFTId(self.ClientSid)
	self.Provider.Send(msg)
}

func (self *Role) SendToSelfMap(msg core.Msg) bool {
	return GS.SendToProvidee(self.MapPvid, msg)
}
