package cmd

import (
	"http-spy-forward/modules"

	"github.com/urfave/cli"
)

var Start = cli.Command{
	Name:        "start",
	Usage:       "sniff local server",
	Description: "startup sniff on local server",
	Action:      modules.Start,
	Flags: []cli.Flag{
		stringFlag("device,i", "", "device name"),
		boolFlag("debug, d", "debug mode"),
		stringFlag("filter,f", "", "setting filters"),
		intFlag("length,l", 1024, "setting snapshot Length"),
		stringFlag("url,u", "", "setting target url"),
	},
}

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func boolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

func intFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}
