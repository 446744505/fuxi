package internal

import "fuxi/core"

type gs struct {
	core.CoreNet
}

func NewGs() *gs {
	g := &gs{}
	return g
}