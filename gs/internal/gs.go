package internal

import (
	"fuxi/core"
	"fuxi/msg"
	"fuxi/providee"
	"github.com/davyxu/golog"
	"strconv"
	"sync"
)

var GS *gs
var Log = golog.New("gs")

type gs struct {
	core.NetControlImpl

	Pvid int32
	roles sync.Map
}

func NewGs() *gs {
	GS = &gs{}
	pvid, _ := strconv.Atoi(core.Args.Get("pvid"))
	GS.Pvid = int32(pvid)
	p := providee.NewProvidee(GS.Pvid, core.ServerGs)
	p.SetPoolCapacity(100)
	p.SetEventHandler(&gsEventHandler{})
	p.SetOnProvideeUpdate(OnProvideeUpdate)
	p.OnClientBroken = OnClientBroken
	GS.AddService(p)
	return GS
}

func OnProvideeUpdate(isRemove bool, meta *core.ProvideeMeta) {
	if core.ServerMap == meta.ServerName {
		if isRemove {
			MapMgr.RemoveMap(meta)
		} else {
			MapMgr.AddMap(meta)
		}
	}
}

func OnClientBroken(clientSid int64) {
	var roleId int64 = 0
	GS.roles.Range(func(_, value interface{}) bool {
		role := value.(*Role)
		if role.ClientSid == clientSid {
			roleId = role.RoleId
			role.ExitGame()
			return false
		}
		return true
	})
	if roleId > 0 {
		GS.roles.Delete(roleId)
	}
}

func (self *gs) OnRoleEnter(p *msg.LEnterGame) {
	role := &Role{}
	role.RoleId = p.RoleId
	role.ClientSid = p.ClientSid
	role.Provider = p.Session()
	self.roles.Store(p.RoleId, role)

	role.EnterGame()
}

func (self *gs) SendToClient(roleId int64, msg core.Msg) {
	role, ok := self.roles.Load(roleId)
	if !ok {
		Log.Errorf("not client role %d", roleId)
		return
	}
	role.(*Role).Send(msg)
	return
}

func (self *gs) SendToProvidee(pvid int32, msg core.Msg) bool {
	return providee.Providee.SendToProvidee(pvid, msg)
}

func (self *gs) OnProviderBroken(session core.Session) {
	self.roles.Range(func(key, value interface{}) bool {
		roleId := key.(int64)
		role := value.(*Role)
		if role.Provider == session {
			role.ExitGame()
			self.roles.Delete(roleId)
		}
		return true
	})
}