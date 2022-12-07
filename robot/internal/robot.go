package internal

import "fuxi/core"

type robot struct {
	core.NetControlImpl
	service core.CoreService
}

func NewRobot() *robot {
	r := &robot{}
	r.AddService(&r.service)
	r.service.SetEventHandler(&robotEventHandler{})
	core.ServiceAddPort(&r.service, core.NewConnector("127.0.0.1", 8080))
	return r
}
