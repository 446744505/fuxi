package internal

import (
	"fmt"
	"fuxi/core"
	"github.com/davyxu/golog"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

var Robot *robot
var Log = golog.New("robot")
var CtxRole = "role"

type robot struct {
	core.NetControlImpl
	service core.CoreService

	linkerLock sync.RWMutex
	linkers map[string]*Linker
}

type linkerWatcher struct {
}

type gsWatcher struct {
}

func NewRobot() *robot {
	Robot = &robot{
		linkers: make(map[string]*Linker),
	}
	Robot.AddService(&Robot.service)
	Robot.service.SetEventHandler(&robotEventHandler{})

	robotNum, _ := strconv.Atoi(core.Args.Get("num"))

	for i := 1; i <= robotNum; i++ {
		NewRole(int64(i)).Start()
	}

	core.ETCD.Watch("linker", &linkerWatcher{})
	return Robot
}

func (self *robot) AddLinker(name, providerName string) {
	self.linkerLock.Lock()
	defer self.linkerLock.Unlock()
	if linker, ok := self.linkers[name]; ok {
		linker.providerName = providerName
	} else {
		self.linkers[name] = &Linker{
			linkerName: name,
			providerName: providerName,
			gsPvids: make(map[int32]bool),
		}
	}
}

func (self *robot) UpdateGs(isRemove bool, providerName string, pvid int32) {
	self.linkerLock.Lock()
	defer self.linkerLock.Unlock()
	for _, linker := range self.linkers {
		if linker.providerName == providerName {
			if isRemove {
				linker.RemoveGs(pvid)
			} else {
				linker.AddGs(pvid)
			}
		}
	}
}

func (self *robot) RandomLinker(gsPvid int32) *Linker {
	self.linkerLock.RLock()
	defer self.linkerLock.RUnlock()
	var linkers []*Linker
	for _, link := range self.linkers {
		if link.HaveGs(gsPvid) {
			linkers = append(linkers, link)
		}
	}
	if len(linkers) == 0 {
		return nil
	}

	return linkers[rand.Intn(len(linkers))]
}

func (self *linkerWatcher) OnAdd(key, val string) {
	arr := strings.Split(key, "/") // key = linker/linkerurl/providerurl
	linkerName := arr[1]
	providerName := arr[2]
	Robot.AddLinker(linkerName, providerName)

	arr = strings.Split(val, ":")
	host := arr[0]
	port, _:= strconv.Atoi(arr[1])
	porter := core.NewConnector(key, host, port)
	if core.ServiceAddPort(&Robot.service, porter) {
		porter.Start()
	}

	core.ETCD.Watch(fmt.Sprintf("providee/%v", providerName), &gsWatcher{})
}

func (self *linkerWatcher) OnDelete(key, val string) {
}

func (self *gsWatcher) OnAdd(key, val string) {
	if val == "gs" {
		arr := strings.Split(key, "/") // key = providee/providerurl/pvid
		pvid, _ := strconv.Atoi(arr[2])
		Robot.UpdateGs(false, arr[1], int32(pvid))
	}
}

func (self *gsWatcher) OnDelete(key, val string) {
	if val == "gs" {
		pvid, _ := strconv.Atoi(val)
		arr := strings.Split(key, "/") // key = providee/providerurl/pvid
		Robot.UpdateGs(true, arr[1], int32(pvid))
	}
}
