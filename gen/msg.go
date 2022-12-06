package msg

import "fuxi/core"

type Dispatch struct {
	core.CoreMsg `binary:"-"`
}

func (self Dispatch) ID() int16 {
	return 1
}

type PDispatch struct {
	core.CoreMsg `binary:"-"`
}

func (self PDispatch) ID() int16 {
	return 2
}

type BindPvid struct {
	core.CoreMsg `binary:"-"`
	PVID int16
}

func (self BindPvid) ID() int16 {
	return 3
}

type UnBindPvid struct {
	core.CoreMsg `binary:"-"`
	PVID int16
}

func (self UnBindPvid) ID() int16 {
	return 4
}