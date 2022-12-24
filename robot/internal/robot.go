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
		NewRole(int64(i), fmt.Sprintf("玩家%v", i)).Start()
	}

	core.ETCD.Watch("linker", &linkerWatcher{})
	core.ETCD.Watch("providee", &gsWatcher{})
	return Robot
}

func (self *robot) AddLinker(name, providerName string) {
	self.linkerLock.Lock()
	defer self.linkerLock.Unlock()
	self.linkers[name] = &Linker{
		linkerName: name,
		providerName: providerName,
		gsPvids: make(map[int32]bool),
	}
}

func (self *robot) AddLinkerSession(session core.Session) {
	self.linkerLock.Lock()
	defer self.linkerLock.Unlock()
	linkerName := session.Port().HostPortString()
	if linker, ok := self.linkers[linkerName]; ok {
		linker.session = session
		Log.Infof("linker add session %v", linkerName)
	} else {
		Log.Errorf("linker add session %v, but linker not exist", linkerName)
	}
}

func (self *robot) RemoveLinkerSession(session core.Session) {
	self.linkerLock.Lock()
	defer self.linkerLock.Unlock()
	linkerName := session.Port().HostPortString()
	if linker, ok := self.linkers[linkerName]; ok {
		linker.session = nil
		Log.Infof(fmt.Sprintf("linker remove session %v", linkerName))
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
		if link.IsAlive() && link.HaveGs(gsPvid) {
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
	Robot.AddLinker(arr[1], arr[2])

	arr = strings.Split(val, ":")
	host := arr[0]
	port, _:= strconv.Atoi(arr[1])
	porter := core.NewConnector(key, host, port)
	if core.ServiceAddPort(&Robot.service, porter) {
		porter.Start()
	}
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
		arr := strings.Split(key, "/") // key = providee/providerurl
		Robot.UpdateGs(true, arr[1], int32(pvid))
	}
}
