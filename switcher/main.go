package main

import (
	"fuxi/core"
	_ "fuxi/msg"
	"fuxi/switcher/linker"
	"fuxi/switcher/provider"
)

func main()  {
	s := &switcher{}
	s.AddService(provider.NewProvider())
	s.AddService(linker.NewLinker())
	s.Start()
	s.Wait()
}

type switcher struct {
	core.NetControlImpl
}