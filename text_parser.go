package ghostdoc

import (
	"github.com/codegangsta/cli"
)

func TextCommand() cli.Command {
	return cli.Command{
		Name:    "text",
		Aliases: []string{"txt"},
		Usage:   "parse text data",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "key, k",
				Usage: "the key used to identify the text segment",
			},
			cli.StringFlag{
				Name:  "pattern, p",
				Usage: "use a pattern file to specify which text segments should be extracted",
			},
		},
		Action: func(c *cli.Context) {
			println("some text stuff")
		},
	}
}
