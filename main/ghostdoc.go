package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/npolar/ghostdoc"
)

func main() {
	InitGhostDoc()
}

func InitGhostDoc() {
	ghostdoc := cli.NewApp()
	ghostdoc.Name = "ghostdoc"
	ghostdoc.Version = "0.0.1"
	ghostdoc.Usage = "Flexible file parser / REST client"
	ghostdoc.Flags = configureFlags()
	ghostdoc.Action = processDocs
	ghostdoc.Commands = defineCommands()
	ghostdoc.Run(os.Args)
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
			Value: "",
			Usage: "Specify keys (before mapping) to exclude from the output",
		},
		cli.StringFlag{
			Name:  "http-verb",
			Value: "POST",
			Usage: "Set the http verb to use [POST|PUT]",
		},
		cli.StringFlag{
			Name:  "key-map, k",
			Usage: "Sets mapping file to use to rename headers/keys. Format {\"oldkey\": \"newkey\"}",
		},
		cli.StringFlag{
			Name:  "merge, m",
			Usage: "Specify additional data to inject into the output",
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
		cli.BoolFlag{
			Name:  "uuid, u",
			Usage: "Injects a namesaced uuid with the 'id' key",
		},
		cli.StringFlag{
			Name:  "wrapper, w",
			Usage: "Define wrapper a wrapper for the payload",
		},
		cli.BoolFlag{
			Name:  "recursive, r",
			Usage: "Recursive read mode. also process sub-dirs",
		},
	}
}

func defineCommands() []cli.Command {
	return []cli.Command{
		ghostdoc.CsvCommand(),
		ghostdoc.JsonCommand(),
		ghostdoc.TextCommand(),
	}

}

func processDocs(c *cli.Context) {
	log.Println("Welocome to", c.App.Name, "See -h for usage info.")
}
