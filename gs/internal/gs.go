package internal

import (
	"fuxi/core"
	"fuxi/providee"
	"github.com/davyxu/golog"
)

var GS *gs
var Log = golog.New("gs")

type gs struct {
	core.NetControlImpl
}

func NewGs() *gs {
	GS = &gs{}
	p := providee.NewProvidee(1, "gs")
	p.SetEventHandler(&gsEventHandler{})
	GS.AddService(p)
	return GS
}