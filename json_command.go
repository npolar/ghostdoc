package ghostdoc

import (
	"github.com/codegangsta/cli"
	"github.com/npolar/ghostdoc/context"
)

// JSONCommand cli.Command for JSON parsing
func JSONCommand() cli.Command {
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
		Action: processJSON,
	}
}

func processJSON(c *cli.Context) {
	jsonStrategy := NewJSONStrategy(context.NewCliContext(c))
	parser := NewParser(jsonStrategy)
	parser.process()
}
