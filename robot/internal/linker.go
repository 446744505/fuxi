package internal

import (
	"fmt"
	"fuxi/core"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

type Linker struct {
	linkerUrl   string
	providerUrl string

	gsLock sync.RWMutex
	gsPvids map[int32]bool
}

func (self *Linker) NewConnect(roleId int64) core.Port {
	arr := strings.Split(self.linkerUrl, ":")
	port, _ := strconv.Atoi(arr[1])
	porter := core.NewConnector(fmt.Sprint(roleId), arr[0], port)
	if core.ServiceAddPort(&Robot.service, porter) {
		porter.Start()
		return porter
	}
	return nil
}

func (self *Linker) HaveGs(pvid int32) bool {
	self.gsLock.RLock()
	defer self.gsLock.RUnlock()
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
	defer self.gsLock.RUnlock()
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
	Log.Infof("robot %v add gs %v", self.linkerUrl, pvid)
}

func (self *Linker) RemoveGs(pvid int32) {
	self.gsLock.Lock()
	defer self.gsLock.Unlock()
	delete(self.gsPvids, pvid)
	Log.Infof("robot %v remove gs %v", self.linkerUrl, pvid)
}
