package internal

import (
	"fuxi/core"
	"math/rand"
	"sync"
)

type Linker struct {
	session core.Session

	linkerName string
	providerName string

	gsLock sync.RWMutex
	gsPvids map[int32]bool
}

func (self *Linker) IsAlive() bool {
	return self.session != nil
}

func (self *Linker) Send(msg core.Msg) bool {
	if self.session == nil {
		return false
	}
	self.session.Send(msg)
	return true
}

func (self *Linker) HaveGs(pvid int32) bool {
	if pvid == 0 {
		return len(self.gsPvids) > 0
	}
	if _, ok := self.gsPvids[pvid]; ok {
		return true
	}

	return false
}

func (self *Linker) RandGs() int32 {
	self.gsLock.RLock()
	self.gsLock.RUnlock()
	var pvids []int32
	for pvid, _ := range self.gsPvids {
		pvids = append(pvids, pvid)
	}
	if len(pvids) == 0 {
		return 0
	}
	return pvids[rand.Intn(len(pvids))]
}

func (self *Linker) AddGs(pvid int32) {
	self.gsLock.Lock()
	defer self.gsLock.Unlock()
	self.gsPvids[pvid] = true
	Log.Infof("linker %v add gs %v", self.linkerName, pvid)
}

func (self *Linker) RemoveGs(pvid int32) {
	self.gsLock.Lock()
	defer self.gsLock.Unlock()
	delete(self.gsPvids, pvid)
	Log.Infof("linker %v remove gs %v", self.linkerName, pvid)
}
