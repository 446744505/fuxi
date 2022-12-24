package internal

import (
	"fuxi/core"
	"fuxi/msg"
	"github.com/davyxu/cellnet"
	"sync"
	"time"
)

type Role struct {
	roleId int64
	name string

	gsPvid int32
	mapPvid int32

	linkerLock sync.RWMutex
	linker *Linker
	session core.Session
}

func NewRole(id int64) *Role {
	return &Role{roleId: id}
}

func (self *Role) Start() {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <- ticker.C:
				self.tryEnterGame()
			}
		}
	}()
}

func (self *Role) OnAddSession(session core.Session) {
	self.session = session
	Log.Debugf("role %v connect success", self.roleId)

	self.linkerLock.RLock()
	gsid := self.linker.RandGs()
	self.linkerLock.RUnlock()
	if gsid == 0 {
		Log.Errorf("role %v no gs can use", self.roleId)
		return
	}

	self.gsPvid = gsid
	Log.Debugf("role %v start enter", self.roleId)
	enter := &msg.CEnterGame{}
	enter.PVID = gsid
	enter.RoleId = self.roleId
	if ok := self.Send(enter); !ok {
		Log.Errorf("role %v enter failed", self.roleId)
	}
}

func (self *Role) OnRemoveSession(_ core.Session) {
	self.session = nil
}

func (self *Role) tryEnterGame() {
	self.linkerLock.Lock()
	defer self.linkerLock.Unlock()
	if self.linker != nil {
		return
	}
	l := Robot.RandomLinker(0)
	if l == nil {
		return
	}

	port := l.NewConnect(self.roleId)
	if port == nil {
		return
	}

	Log.Debugf("role %v start connect", self.roleId)
	ctx := port.Peer().(cellnet.ContextSet)
	ctx.SetContext(CtxRole, self)
	self.linker = l
}

func (self *Role) Send(msg core.Msg) bool {
	if self.session == nil {
		return false
	}
	self.session.Send(msg)
	return true
}

func (self *Role) EnterGame(enter *msg.SEnterGame) {
	self.name = enter.Name
	Log.Infof("role %s enter gs", enter.Name)
}

func (self *Role) EnterMap(enter *msg.SEnterMap) {
	self.mapPvid = enter.Pvid
	Log.Infof("role %s enter map %v", self.name, enter.Pvid)
}