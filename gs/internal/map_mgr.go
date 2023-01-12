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

func (self *mapMgr) RandomMap() int32 {
	var maps []int32
	for pvid := range self.maps {
		maps = append(maps, pvid)
	}
	if len(maps) == 0 {
		return 0
	}

	return maps[rand.Intn(len(maps))]
}

func (self *mapMgr) AddMap(meta *core.ProvideeMeta) {
	if core.ServerMap != meta.ServerName {
		return
	}
	self.maps[meta.Pvid] = &mapInfo{providerUrl: meta.ProviderUrl}
	Log.Infof("add map %v", meta.Pvid)
}

func (self *mapMgr) RemoveMap(meta *core.ProvideeMeta) {
	if core.ServerMap != meta.ServerName {
		return
	}
	delete(self.maps, meta.Pvid)
	Log.Infof("remove map %v", meta.Pvid)
}

func init() {
	MapMgr = &mapMgr{
		maps: make(map[int32]*mapInfo),
	}
}
