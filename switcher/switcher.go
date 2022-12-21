package switcher

import (
	"fuxi/core"
	"github.com/urfave/cli/v2"
)

var Args *cli.Context

type Switcher struct {
	core.NetControlImpl
}
