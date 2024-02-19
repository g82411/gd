package cmd

import (
	"context"
	"github.com/urfave/cli/v2"
)

const (
	appName = "gd"
)

var app = &cli.App{
	Name:                 appName,
	Usage:                "Control google drive file using cli command",
	EnableBashCompletion: true,
	Action: func(c *cli.Context) error {
		args := c.Args()
		if args.Present() {
			cli.ShowCommandHelp(c, args.First())
			return cli.Exit("", 1)
		}
		return cli.ShowAppHelp(c)
	},
}

type Result struct {
	Path     string
	FileName string
	Url      string
	Size     int64
}

func Main(ctx context.Context, args []string) error {
	app.Commands = []*cli.Command{
		NewListCommand(),
	}
	return app.RunContext(ctx, args)
}
