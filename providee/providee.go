package providee

import (
	"fmt"
	"fuxi/net"
)

type Providee struct {
	ID int16
	Name string
}

func NewProvidee(id int16, name string) *Providee {
	return &Providee{ID: id, Name: name}
}

func (self *Providee) Start() {
	factory := &ProvideeFactory{}
	netConf := net.NewNetConf(self.Name)
	provideeConf := netConf.NewService(fmt.Sprintf("providee-%s", self.Name)).BuildFactory(factory)
	//todo 从etcd拿所有provider
	provideeConf.NewPort(net.Connector).BuildHost("127.0.0.1").BuildPort(8088)
	ops := net.NewNetOptions(netConf).BuildFactory(factory)
	net := net.FXNet(ops)
	net.Start()
	net.Wait()
	net.Stop()
}
