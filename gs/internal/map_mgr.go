package internal

import (
	"fuxi/core"
	"math/rand"
	"strconv"
	"strings"
)

var MapMgr *mapMgr

type mapInfo struct {
	providerName string
}

type mapMgr struct {
	maps map[int32]*mapInfo //key=pvid
}

func (self *mapMgr) Init() {
	core.ETCD.Watch("providee", self)
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
		arr := strings.Split(key, "/") // key = providee/providerurl/pvid
		pvid, _ := strconv.Atoi(arr[2])
		self.maps[int32(pvid)] = &mapInfo{providerName: arr[1]}
		Log.Infof("add map %v", pvid)
	}
}

func (self *mapMgr) OnDelete(key, val string) {
	if val == "map" {
		arr := strings.Split(key, "/") // key = providee/providerurl/pvid
		pvid, _ := strconv.Atoi(arr[2])
		delete(self.maps, int32(pvid))
		Log.Infof("remove map %v", pvid)
	}
}

func init() {
	MapMgr = &mapMgr{
		maps: make(map[int32]*mapInfo),
	}
}
