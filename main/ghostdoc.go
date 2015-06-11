package main

import (
	"os"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/npolar/ghostdoc"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := initGhostDoc().Run(os.Args)
	if err != nil {
		log.Error(err.Error())
	}
}

func initGhostDoc() *cli.App {
	ghostdoc := cli.NewApp()
	ghostdoc.Name = "ghostdoc"
	ghostdoc.Version = "0.0.1"
	ghostdoc.Usage = "Flexible file parser / REST client"
	ghostdoc.Flags = configureFlags()
	ghostdoc.Action = processDocs
	ghostdoc.Commands = defineCommands()
	return ghostdoc
}

func configureFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "address, a",
			Usage: "Set url to write to",
		},
		cli.IntFlag{
			Name:  "concurrency, c",
			Value: 2,
			Usage: "Specify the number of concurrent operations",
		},
		cli.StringFlag{
			Name:  "exclude, e",
			Usage: "Specify keys (before mapping) to exclude in the output",
		},
		cli.StringFlag{
			Name:  "filename, f",
			Usage: "Set filename to use in name-pattern when piping data via stdin",
		},
		cli.StringFlag{
			Name:  "include, i",
			Usage: "Specify keys (before mapping) to include in the output",
		},
		cli.StringFlag{
			Name:  "js, j",
			Usage: "Run javascript map functions on the data",
		},
		cli.StringFlag{
			Name:  "http-verb",
			Value: "POST",
			Usage: "Set the http verb to use [POST|PUT]",
		},
		cli.StringFlag{
			Name:  "key-map, k",
			Usage: "Sets mapping file to use to rename headers/keys. JSON Format {\"oldkey\": \"newkey\"}",
		},
		cli.StringFlag{
			Name:  "merge, m",
			Usage: "Specify additional JSON data to inject into the output.",
		},
		cli.StringFlag{
			Name:  "name-pattern, n",
			Usage: "Set pattern file to extract filename info and inject it into the result",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Set dir output dir. Files will get uuid as name",
		},
		cli.StringFlag{
			Name:  "payload-key, p",
			Value: "data",
			Usage: "Specify the key to use for the payload when wrapping",
		},
		cli.StringFlag{
			Name:  "log-file, lf",
			Usage: "Log to file instead of stdout",
		},
		cli.StringFlag{
			Name:  "log-level, ll",
			Value: "DEBUG",
			Usage: "Log level DEBUG|INFO|ERROR|OFF",
		},
		cli.StringFlag{
			Name:  "log-mail, lm",
			Usage: "Forward log errors to email",
		},
		cli.BoolFlag{
			Name:  "uuid, u",
			Usage: "Injects a namesaced uuid with the 'id' key",
		},
		cli.StringFlag{
			Name:  "uuid-include, ui",
			Usage: "Injects a namesaced uuid with the 'id' key based on a set of keys",
		},
		cli.StringFlag{
			Name:  "uuid-key, uk",
			Value: "id",
			Usage: "uuid key name",
		},
		cli.StringFlag{
			Name:  "wrapper, w",
			Usage: "Define JSON wrapper a wrapper for the payload",
		},
		cli.BoolFlag{
			Name:  "recursive, r",
			Usage: "Recursive read mode. also process sub-dirs",
		},
		cli.StringFlag{
			Name:  "schema, s",
			Usage: "Reference to a JSON Schema to validate json output against",
		},
	}
}

func defineCommands() []cli.Command {
	return []cli.Command{
		ghostdoc.CsvCommand(),
		ghostdoc.JSONCommand(),
		ghostdoc.TextCommand(),
	}

}

func processDocs(c *cli.Context) {
	cli.ShowAppHelp(c)
}
