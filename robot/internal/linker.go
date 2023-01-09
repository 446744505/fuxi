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

	gsPvids sync.Map
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
	if pvid == 0 {
		have := false
		self.gsPvids.Range(func(key, value interface{}) bool {
			have = true
			return false
		})
		return have
	}
	if _, ok := self.gsPvids.Load(pvid); ok {
		return true
	}

	return false
}

func (self *Linker) RandGs() int32 {
	var pvids []int32
	self.gsPvids.Range(func(key, _ interface{}) bool {
		pvid := key.(int32)
		pvids = append(pvids, pvid)
		return true
	})
	if len(pvids) == 0 {
		return 0
	}
	return pvids[rand.Intn(len(pvids))]
}

func (self *Linker) AddGs(pvid int32) {
	self.gsPvids.Store(pvid, struct{}{})
	Log.Infof("robot %v add gs %v", self.linkerUrl, pvid)
}

func (self *Linker) RemoveGs(pvid int32) {
	self.gsPvids.Delete(pvid)
	Log.Infof("robot %v remove gs %v", self.linkerUrl, pvid)
}
