package main

import (
	_ "fuxi/gen"
	"fuxi/net"
)

func main()  {
	netConf := net.NewNetConf("switcher")
	linkerConf := netConf.NewService("linker")
	linkerConf.NewPort(net.Acceptor).BuildHost("127.0.0.1").BuildPort(8080)
	providerConf := netConf.NewService("provider")
	providerConf.NewPort(net.Acceptor).BuildHost("127.0.0.1").BuildPort(8088)
	ops := net.NewNetOptions(netConf).BuildFactory(&SwitcherFactory{})
	net := net.FXNet(ops)
	net.Start()
	net.Wait()
	net.Stop()
}

type SwitcherFactory struct {
	net.SampleNetFactory
}
