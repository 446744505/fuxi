package linker

import "fuxi/core"

var Linker *linker

type linker struct {
	core.CoreService
}

func NewLinker() *linker {
	Linker = &linker{}
	Linker.SetName("linker")
	Linker.NewPort(func() core.Port {
		return core.NewAcceptor("127.0.0.1", 8080)
	})
	return Linker
}
