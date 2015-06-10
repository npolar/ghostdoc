package ghostdoc

import (
	"github.com/codegangsta/cli"
	"github.com/npolar/ghostdoc/context"
)

// TextCommand cli.Command for parsing text
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
		Action: processText,
	}
}

func processText(c *cli.Context) {
	textStrategy := NewTextStrategy(context.NewCliContext(c))
	parser := NewParser(textStrategy)
	parser.process()
}
