package internal

import (
	"fuxi/core"
	"math/rand"
)

var MapMgr *mapMgr

type mapInfo struct {
	providerUrl string
}

type mapMgr struct {
	maps map[int32]*mapInfo //key=pvid
}

func (self *mapMgr) Init() {
	core.ETCD.Watch(core.NodeNameProvidee, self)
}

func (self *mapMgr) RandomMap() int32 {
	var maps []int32
	for pvid, _ := range self.maps {
		maps = append(maps, pvid)
	}
	if len(maps) == 0 {
		return 0
	}

	return maps[rand.Intn(len(maps))]
}

func (self *mapMgr) OnAdd(key, val string) {
	if val == "map" {
		meta := &core.ProvideeMeta{}
		meta.ValueOf(key)
		self.maps[meta.Pvid] = &mapInfo{providerUrl: meta.ProviderUrl}
		Log.Infof("add map %v", meta.Pvid)
	}
}

func (self *mapMgr) OnDelete(key, val string) {
	if val == "map" {
		meta := &core.ProvideeMeta{}
		meta.ValueOf(key)
		delete(self.maps, meta.Pvid)
		Log.Infof("remove map %v", meta.Pvid)
	}
}

func init() {
	MapMgr = &mapMgr{
		maps: make(map[int32]*mapInfo),
	}
}
