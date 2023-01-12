package main

import (
	"fuxi/core"
	"fuxi/switcher/linker"
	"fuxi/switcher/provider"
	"github.com/davyxu/golog"
	"strings"
)

func main() {
	args := core.CreateArgs("Switcher", "the fuxi switcher")
	args.Flag(core.ArgEtcd, "127.0.0.1:2379", "the etcd url")
	args.Flag("linker", "127.0.0.1:8080", "the linker url")
	args.Flag("provider", "127.0.0.1:8088", "the provider url")
	args.Flag(core.ArgLogLevel, "debug", "the log level")
	args.Flag(core.ArgPPort, "0", "the pprof port")

	if err := args.Run(); err != nil {
		core.Log.Errorf("%v", err)
		return
	}

	golog.SetLevelByString(".", core.Args.Get(core.ArgLogLevel))
	url := core.Args.Get(core.ArgEtcd)
	core.InitEtcd(strings.Split(url, ","))

	s := &switcher{}
	s.SetPProfPort(core.Args.GetInt(core.ArgPPort))
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
