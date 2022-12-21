package main

import (
	"fuxi/core"
	"fuxi/switcher"
	"fuxi/switcher/linker"
	"fuxi/switcher/provider"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

func main()  {
	app := &cli.App{
		Name: "Switcher",
		Usage: "the fuxi switcher",
		Action: func(c *cli.Context) error {
			switcher.Args = c
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "etcd",
				Value: "127.0.0.1:2379",
				Usage: "the etcd url",
			},
			&cli.StringFlag{
				Name: "linker",
				Value: "127.0.0.1:8080",
				Usage: "the linker url",
			},
			&cli.StringFlag{
				Name: "provider",
				Value: "127.0.0.1:8088",
				Usage: "the provider url",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		core.Log.Errorf("%v", err)
		return
	}

	url := switcher.Args.String("etcd")
	core.InitEtcd(strings.Split(url, ","))

	s := &switcher.Switcher{}
	s.AddService(provider.NewProvider())
	s.AddService(linker.NewLinker())
	s.Start()
	s.Wait()
}