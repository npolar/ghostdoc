package ghostdoc

import (
	"github.com/codegangsta/cli"
	"github.com/npolar/ghostdoc/context"
)

// CsvCommand specifies the cli interface for the csv parser
func CsvCommand() cli.Command {
	return cli.Command{
		Name:  "csv",
		Usage: "Parse delimiter separated value files (csv, tsv, etc...)",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "delimiter, d",
				Value: ",",
				Usage: "Set the delimiter char.",
			},
			cli.StringFlag{
				Name:  "comment, c",
				Value: "#",
				Usage: "Set the comment char.",
			},
			cli.StringFlag{
				Name:  "header, hd",
				Usage: "Configure data header. If not set the first data line will be used.",
			},
			cli.IntFlag{
				Name:  "skip, s",
				Usage: "Specify the number of lines to skip before parsing. [NOTE] Blank lines are ignored by the parser and should not be skipped.",
			},
		},
		Action: processCsv,
	}
}

func processCsv(c *cli.Context) {
	csvStrategy := NewCsvStrategy(context.NewCliContext(c))
	parser := NewParser(csvStrategy)
	parser.process()
}
