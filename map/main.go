package main

import (
	"fuxi/core"
	"fuxi/map/internal"
	"github.com/davyxu/golog"
	"strings"
)

func main() {
	args := core.CreateArgs("map", "the fuxi map")
	args.Flag(core.ArgEtcd, "127.0.0.1:2379", "the etcd url")
	args.Flag(core.ArgPvid, "0", "the providee id")
	args.Flag(core.ArgLogLevel, "debug", "the log level")
	args.Flag(core.ArgPPort, "0", "the pprof port")

	if err := args.Run(); err != nil {
		internal.Log.Errorf("%v", err)
		return
	}

	golog.SetLevelByString(".", core.Args.Get(core.ArgLogLevel))
	url := core.Args.Get(core.ArgEtcd)
	core.InitEtcd(strings.Split(url, ","))

	m := internal.NewMap()
	m.Start()
	internal.Log.Infof("map server started")
	m.Wait()
	core.StopEtcd()
	internal.Log.Infof("map server stoped")
}
