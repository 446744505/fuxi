package internal

import (
	"fuxi/core"
	"github.com/davyxu/golog"
	"strconv"
	"strings"
)

var Log = golog.New("robot")

type robot struct {
	core.NetControlImpl
	service core.CoreService
}

func NewRobot() *robot {
	r := &robot{}
	r.AddService(&r.service)
	r.service.SetEventHandler(&robotEventHandler{})

	core.ETCD.Watch("linker", r)

	return r
}

func (self *robot) OnAdd(key, val string) {
	arr := strings.Split(val, ":")
	host := arr[0]
	port, _:= strconv.Atoi(arr[1])
	porter := core.NewConnector(key, host, port)
	core.ServiceAddPort(&self.service, porter)
	porter.Start()
}

func (self *robot) OnDelete(key, val string) {

}
