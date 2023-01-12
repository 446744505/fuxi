package main

import (
	"fuxi/core"
	"fuxi/robot/internal"
	"github.com/davyxu/golog"
	"strings"
)

func main() {
	args := core.CreateArgs("robot", "the fuxi robot")
	args.Flag("etcd", "127.0.0.1:2379", "the etcd url")
	args.Flag("num", "1", "the number of robot")
	args.Flag("log_level", "debug", "the log level")

	if err := args.Run(); err != nil {
		internal.Log.Errorf("%v", err)
		return
	}

	golog.SetLevelByString(".", core.Args.Get("log_level"))
	url := core.Args.Get("etcd")
	core.InitEtcd(strings.Split(url, ","))

	r := internal.NewRobot()
	r.Start()
	internal.Log.Infof("robot client started")
	r.Wait()
	core.StopEtcd()
	internal.Log.Infof("robot client stoped")
}

