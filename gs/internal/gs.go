package internal

import (
	"fuxi/core"
	"fuxi/providee"
)

var GS *gs

type gs struct {
	core.NetControlImpl
}

func NewGs() *gs {
	GS := &gs{}
	p := providee.NewProvidee(1, "gs")
	p.SetEventHandler(&gsEventHandler{})
	GS.AddService(p)
	return GS
}