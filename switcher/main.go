package main

import (
	"fuxi/core"
	_ "fuxi/gen"
	"fuxi/switcher/linker"
	"fuxi/switcher/provider"
)

func main()  {
	s := &Switcher{}
	s.NewService(func() core.Service {
		return linker.NewLinker()
	})
	s.NewService(func() core.Service {
		return provider.NewProvider()
	})
	s.Start()
	s.Wait()
}

type Switcher struct {
	core.CoreNet
}