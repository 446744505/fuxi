package main

import (
	"fuxi/core"
	"fuxi/switcher/linker"
	"fuxi/switcher/provider"
	"github.com/davyxu/golog"
	"strings"
)

func main()  {
	golog.SetLevelByString(".", "info")
	args := core.CreateArgs("Switcher", "the fuxi switcher")
	args.Flag("etcd", "127.0.0.1:2379", "the etcd url")
	args.Flag("linker", "127.0.0.1:8080", "the linker url")
	args.Flag("provider", "127.0.0.1:8088", "the provider url")

	if err := args.Run(); err != nil {
		core.Log.Errorf("%v", err)
		return
	}

	url := core.Args.Get("etcd")
	core.InitEtcd(strings.Split(url, ","))

	s := &switcher{}
	s.AddService(provider.NewProvider())
	s.AddService(linker.NewLinker())
	s.Start()
	core.Log.Infof("switcher server started")
	s.Wait()
	core.StopEtcd()
	core.Log.Infof("switcher server stoped")
}

type switcher struct {
	core.NetControlImpl
}