package internal

import (
	"fuxi/core"
	"fuxi/providee"
	"github.com/davyxu/golog"
	"strconv"
)

var Map *mmap
var Log = golog.New("map")

type mmap struct {
	core.NetControlImpl
}

func NewMap() *mmap {
	Map = &mmap{}
	pvid, _ := strconv.Atoi(core.Args.Get("pvid"))
	p := providee.NewProvidee(int32(pvid), "map")
	p.SetEventHandler(&mapEventHandler{})
	Map.AddService(p)
	return Map
}
