package msg

type Dispatch struct {

}

func (self Dispatch) ID() int16 {
	return 1
}

type PDispatch struct {

}

func (self PDispatch) ID() int16 {
	return 2
}

type BindPvid struct {
	PVID int16
}

func (self BindPvid) ID() int16 {
	return 3
}

type UnBindPvid struct {
	PVID int16
}

func (self UnBindPvid) ID() int16 {
	return 4
}