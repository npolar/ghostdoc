package main

import (
	"flag"
	"github.com/codegangsta/cli"
)

func main() {
	InitGhostDoc()
}

func InitGhostDoc() {
	ghostdoc := cli.NewApp()
	ghostdoc.Name = "ghostdoc"
	ghostdoc.Version = "0.0.1"
	ghostdoc.Usage = "Client used posting data to REST apis"
	ghostdoc.Flags = ConfigureFlags()
}

func ConfigureFlags() []cli.Flag {
	return []cli.Flag{
		cli.IntFlag{
			Name:  "concurrency, c",
			Value: 2,
			Usage: "Specify the number of concurrent operations",
		},
		cli.StringFlag{
			Name:  "delimiter, d",
			Value: "i",
			Usage: "Set the demlimiter for working with csv data",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "json",
			Usage: "Set the input format to watch for [json|csv|txt].",
		},
		cli.StringFlag{
			Name:  "http-verb",
			Value: "POST",
			Usage: "Set the http verb to use [POST|PUT]",
		},
		cli.StringFlag{
			Name:  "key, k",
			Value: "data",
			Usage: "Specify the key to use for the data when wrapping",
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
			Name:  "uuid, u",
			Usage: "Set namespaced uuid generation. uuid's will be injected with the 'id' key [full|wrap].",
		},
		cli.StringFlag{
			Name:  "wrapper, w",
			Usage: "Define wrapper content for the data",
		},
		cli.StringFlag{
			Name:  "remap, r",
			Usage: "Sets mapping file to use to rename keys. Format {\"oldkey\": \"newkey\"}",
		},
	}
}
