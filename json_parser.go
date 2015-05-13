package ghostdoc

import (
	"github.com/codegangsta/cli"
)

func JsonCommand() cli.Command {
	return cli.Command{
		Name:  "json",
		Usage: "parse json files",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "bulk, b",
				Usage: "enable bulk processing",
			},
			cli.StringFlag{
				Name:  "map, m",
				Usage: "rename keys. Format: '{\"oldKey\": \"newKey\"}'",
			},
		},
		Action: func(c *cli.Context) {
			println("hello json")
		},
	}
}
