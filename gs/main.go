package main

import (
	"fuxi/core"
	"fuxi/gs/internal"
)

func main() {
	core.InitEtcd([]string{"127.0.0.1:2379"})
	g := internal.NewGs()
	g.Start()
	g.Wait()
}