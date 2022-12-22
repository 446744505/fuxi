package internal

import (
	"math/rand"
	"strconv"
)

var MapMgr *mapMgr

type mapInfo struct {

}

type mapMgr struct {
	maps map[int32]*mapInfo //key=pvid
}

func (self *mapMgr) Init() {
	//core.ETCD.Watch("map", self)
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
	pvid, _ := strconv.Atoi(val)
	self.maps[int32(pvid)] = &mapInfo{}
	Log.Infof("add map %v", pvid)
}

func (self *mapMgr) OnDelete(key, val string) {
	pvid, _ := strconv.Atoi(val)
	delete(self.maps, int32(pvid))
	Log.Infof("remove map %v", pvid)
}

func init() {
	MapMgr = &mapMgr{
		maps: make(map[int32]*mapInfo),
	}
}
