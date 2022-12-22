package internal

import (
	"fuxi/core"
	"fuxi/msg"
	"time"
)

type Role struct {
	roleId int64
	name string

	linker *Linker
}

func NewRole(id int64, name string) *Role {
	return &Role{roleId: id, name: name}
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

func (self *Role) tryEnterGame() {
	if self.linker != nil {
		return
	}
	l := Robot.RandomLinker(0)
	if l == nil {
		return
	}
	gsid := l.RandGs()
	if gsid == 0 {
		return
	}

	self.linker = l
	enter := &msg.CEnterGame{}
	enter.PVID = gsid
	enter.RoleId = self.roleId
	self.Send(enter)
}

func (self *Role) Send(msg core.Msg) bool {
	if self.linker == nil {
		return false
	}
	return self.linker.Send(msg)
}