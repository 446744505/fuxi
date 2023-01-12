package main

import (
	"fuxi/core"
	"fuxi/robot/internal"
	"github.com/davyxu/golog"
	"strings"
)

func main() {
	args := core.CreateArgs("robot", "the fuxi robot")
	args.Flag(core.ArgEtcd, "127.0.0.1:2379", "the etcd url")
	args.Flag("num", "1", "the number of robot")
	args.Flag(core.ArgLogLevel, "debug", "the log level")
	args.Flag(core.ArgPPort, "0", "the pprof port")

	if err := args.Run(); err != nil {
		internal.Log.Errorf("%v", err)
		return
	}

	golog.SetLevelByString(".", core.Args.Get(core.ArgLogLevel))
	url := core.Args.Get(core.ArgEtcd)
	core.InitEtcd(strings.Split(url, ","))

	r := internal.NewRobot()
	r.SetPProfPort(core.Args.GetInt(core.ArgPPort))
	r.Start()
	internal.Log.Infof("robot client started")
	r.Wait()
	core.StopEtcd()
	internal.Log.Infof("robot client stoped")
}
