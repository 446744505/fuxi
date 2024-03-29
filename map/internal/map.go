package internal

import (
	"fuxi/core"
	"fuxi/msg"
	"fuxi/providee"
	"github.com/davyxu/golog"
	"sync"
)

var Map *mmap
var Log = golog.New("map")

type mmap struct {
	core.NetControlImpl

	pvid  int32
	roles sync.Map
}

func NewMap() *mmap {
	Map = &mmap{}
	pvid := core.Args.GetInt(core.ArgPvid)
	Map.pvid = int32(pvid)
	Map.SetPProfPort(core.Args.GetInt(core.ArgPPort))
	p := providee.NewProvidee(int32(pvid), core.ServerMap)
	p.SetEventHandler(&mapEventHandler{})
	Map.AddService(p)
	return Map
}

func (self *mmap) OnRoleEnter(enter *msg.GEnterMap) {
	role := &Role{}
	self.roles.Store(enter.RoleId, role)
	role.OnEnterMap(enter)
}

func (self *mmap) OnRoleExit(roleId int64) {
	if role := self.GetRole(roleId); role != nil {
		self.roles.Delete(roleId)
		role.OnExitMap()
	}
}

func (self *mmap) GetRole(roleId int64) *Role {
	if role, ok := self.roles.Load(roleId); ok {
		return role.(*Role)
	}
	return nil
}

func (self *mmap) SendToClient(providerName string, clientSid int64, msg core.Msg) bool {
	if provider := providee.Providee.GetProvider(providerName); provider == nil {
		return false
	} else {
		msg.SetFTId(clientSid)
		msg.SetToType(core.MsgToClient)
		provider.Send(msg)
	}
	return true
}

func (self *mmap) SendToProvidee(pvid int32, msg core.Msg) bool {
	return providee.Providee.SendToProvidee(pvid, msg)
}
