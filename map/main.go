package main

import (
	"fuxi/core"
	"fuxi/map/internal"
	"github.com/davyxu/golog"
	"strings"
)

func main() {
	args := core.CreateArgs("map", "the fuxi map")
	args.Flag("etcd", "127.0.0.1:2379", "the etcd url")
	args.Flag("pvid", "0", "the providee id")
	args.Flag("log_level", "debug", "the log level")

	if err := args.Run(); err != nil {
		internal.Log.Errorf("%v", err)
		return
	}

	golog.SetLevelByString(".", core.Args.Get("log_level"))
	url := core.Args.Get("etcd")
	core.InitEtcd(strings.Split(url, ","))

	m := internal.NewMap()
	m.Start()
	internal.Log.Infof("map server started")
	m.Wait()
	core.StopEtcd()
	internal.Log.Infof("map server stoped")
}
