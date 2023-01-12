package core

import (
	"github.com/urfave/cli/v2"
	"os"
)

var Args *args

type args struct {
	app *cli.App
	ctx *cli.Context
}

func CreateArgs(name, usage string) *args {
	Args = &args{}
	Args.app = &cli.App{
		Name:  name,
		Usage: usage,
		Action: func(context *cli.Context) error {
			Args.ctx = context
			return nil
		},
	}
	return Args
}

func (self *args) Flag(name, value, usage string) *args {
	self.app.Flags = append(self.app.Flags, &cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	})
	return self
}

func (self *args) Get(key string) string {
	return self.ctx.String(key)
}

func (self *args) Run() error {
	if err := self.app.Run(os.Args); err != nil {
		Log.Errorf("%v", err)
		return err
	}
	return nil
}
