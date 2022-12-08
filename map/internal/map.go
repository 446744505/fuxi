package internal

import (
	"fuxi/core"
	"fuxi/providee"
)

var Map *mmap

type mmap struct {
	core.NetControlImpl
}

func NewMap() *mmap {
	Map = &mmap{}
	p := providee.NewProvidee(2, "map")
	p.SetEventHandler(&mapEventHandler{})
	Map.AddService(p)
	return Map
}
