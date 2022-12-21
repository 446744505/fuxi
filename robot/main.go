package main

import (
	"fuxi/core"
	"fuxi/robot/internal"
	"strings"
)

func main() {
	args := core.CreateArgs("robot", "the fuxi robot")
	args.Flag("etcd", "127.0.0.1:2379", "the etcd url")

	if err := args.Run(); err != nil {
		internal.Log.Errorf("%v", err)
		return
	}

	url := core.Args.Get("etcd")
	core.InitEtcd(strings.Split(url, ","))

	r := internal.NewRobot()
	r.Start()
	r.Wait()
}

