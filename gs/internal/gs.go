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
	roleLock sync.RWMutex
	roles map[int64]*NetRole
}

func NewGs() *gs {
	GS = &gs{
		roles: make(map[int64]*NetRole),
	}
	pvid, _ := strconv.Atoi(core.Args.Get("pvid"))
	GS.Pvid = int32(pvid)
	p := providee.NewProvidee(GS.Pvid, core.ServerGs)
	p.SetEventHandler(&gsEventHandler{})
	p.SetOnProvideeUpdate(OnProvideeUpdate)
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

func (self *gs) OnRoleEnter(p *msg.LEnterGame) {
	role := &NetRole{}
	role.RoleId = p.RoleId
	role.ClientSid = p.ClientSid
	role.Provider = p.Session()

	self.roleLock.Lock()
	self.roles[p.RoleId] = role
	self.roleLock.Unlock()

	role.EnterGame()
}

func (self *gs) SendToClient(roleId int64, msg core.Msg) {
	self.roleLock.RLock()
	role, ok := self.roles[roleId]
	self.roleLock.RUnlock()
	if !ok {
		Log.Errorf("not client role %d", roleId)
		return
	}
	role.Send(msg)
	return
}

func (self *gs) SendToProvidee(pvid int32, msg core.Msg) bool {
	return providee.Providee.SendToProvidee(pvid, msg)
}

func (self *gs) OnProviderBroken(session core.Session) {
	self.roleLock.Lock()
	defer self.roleLock.Unlock()
	for roleId, role := range self.roles {
		if role.Provider == session {
			role.ExitGame()
			delete(self.roles, roleId)
		}
	}
}