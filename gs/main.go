package main

import (
	"fuxi/core"
	"fuxi/gs/internal"
	"strings"
)

func main() {
	args := core.CreateArgs("gs", "the fuxi gs")
	args.Flag("etcd", "127.0.0.1:2379", "the etcd url")
	args.Flag("pvid", "0", "the providee id")

	if err := args.Run(); err != nil {
		internal.Log.Errorf("%v", err)
		return
	}

	url := core.Args.Get("etcd")
	core.InitEtcd(strings.Split(url, ","))

	g := internal.NewGs()
	g.Start()
	g.Wait()
	core.StopEtcd()
}