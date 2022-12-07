package internal

import "fuxi/core"

type robot struct {
	core.CoreNet
}

type robotService struct {
	core.CoreService
}

func NewRobot() *robot {
	r := &robot{}
	s := r.NewService(func() core.Service {
		return &robotService{}
	})
	s.NewPort(func() core.Port {
		return core.NewConnector("127.0.0.1", 8080)
	})
	return r
}
