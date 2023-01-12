package internal

import (
	"fmt"
	"fuxi/core"
	"github.com/davyxu/golog"
	"math/rand"
	"sync"
)

var Robot *robot
var Log = golog.New("robot")
var CtxTypeRole = "role"

type robot struct {
	core.NetControlImpl
	service core.CoreService

	linkers sync.Map
}

type linkerWatcher struct {
}

type gsWatcher struct {
}

func NewRobot() *robot {
	Robot = &robot{}
	Robot.AddService(&Robot.service)
	Robot.service.SetEventHandler(&robotEventHandler{})

	robotNum := core.Args.GetInt("num")

	for i := 1; i <= robotNum; i++ {
		NewRole(int64(i)).Start()
	}

	core.ETCD.Watch(core.NodeNameLinker, &linkerWatcher{})
	return Robot
}

func (self *robot) AddLinker(linkerUrl, providerUrl string) {
	if linker, ok := self.linkers.Load(linkerUrl); ok {
		linker.(*Linker).providerUrl = providerUrl
	} else {
		linker := &Linker{
			linkerUrl:   linkerUrl,
			providerUrl: providerUrl,
		}
		self.linkers.Store(linkerUrl, linker)
	}
}

func (self *robot) UpdateGs(isRemove bool, providerUrl string, pvid int32) {
	self.linkers.Range(func(_, value interface{}) bool {
		linker := value.(*Linker)
		if linker.providerUrl == providerUrl {
			if isRemove {
				linker.RemoveGs(pvid)
			} else {
				linker.AddGs(pvid)
			}
		}
		return true
	})
}

func (self *robot) RandomLinker(gsPvid int32) *Linker {
	var linkers []*Linker
	self.linkers.Range(func(_, value interface{}) bool {
		linker := value.(*Linker)
		if linker.HaveGs(gsPvid) {
			linkers = append(linkers, linker)
		}
		return true
	})

	if len(linkers) == 0 {
		return nil
	}

	return linkers[rand.Intn(len(linkers))]
}

func (self *linkerWatcher) OnAdd(key, val string) {
	meta := &core.SwitcherMeta{}
	meta.ValueOf(key)
	Robot.AddLinker(meta.LinkerUrl, meta.ProviderUrl)

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
