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
var CtxTypeRole = "role"

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

	core.ETCD.Watch(core.NodeNameLinker, &linkerWatcher{})
	return Robot
}

func (self *robot) AddLinker(linkerUrl, providerUrl string) {
	self.linkerLock.Lock()
	defer self.linkerLock.Unlock()
	if linker, ok := self.linkers[linkerUrl]; ok {
		linker.providerUrl = providerUrl
	} else {
		self.linkers[linkerUrl] = &Linker{
			linkerUrl:   linkerUrl,
			providerUrl: providerUrl,
			gsPvids:     make(map[int32]bool),
		}
	}
}

func (self *robot) UpdateGs(isRemove bool, providerUrl string, pvid int32) {
	self.linkerLock.Lock()
	defer self.linkerLock.Unlock()
	for _, linker := range self.linkers {
		if linker.providerUrl == providerUrl {
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
	meta := &core.SwitcherMeta{}
	meta.ValueOf(key)
	Robot.AddLinker(meta.LinkerUrl, meta.ProviderUrl)

	arr := strings.Split(val, ":")
	host := arr[0]
	port, _:= strconv.Atoi(arr[1])
	porter := core.NewConnector(key, host, port)
	if core.ServiceAddPort(&Robot.service, porter) {
		porter.Start()
	}

	core.ETCD.Watch(fmt.Sprintf("%s/%s", core.NodeNameProvidee, meta.ProviderUrl), &gsWatcher{})
}

func (self *linkerWatcher) OnDelete(key, val string) {
}

func (self *gsWatcher) OnAdd(key, val string) {
	if val == core.ServerGs {
		meta := &core.ProvideeMeta{}
		meta.ValueOf(key, val)
		Robot.UpdateGs(false, meta.ProviderUrl, meta.Pvid)
	}
}

func (self *gsWatcher) OnDelete(key, val string) {
	if val == core.ServerGs {
		meta := &core.ProvideeMeta{}
		meta.ValueOf(key, val)
		Robot.UpdateGs(true, meta.ProviderUrl, meta.Pvid)
	}
}
