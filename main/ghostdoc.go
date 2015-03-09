package main

import (
	"github.com/codegangsta/cli"
	"os"
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
	ghostdoc.Action = ProcessDocs
	ghostdoc.Run(os.Args)
}

func ConfigureFlags() []cli.Flag {
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
			Name:  "delimiter, d",
			Value: ",",
			Usage: "Set the demlimiter for working with csv data",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "json",
			Usage: "Set the input format [json|csv|txt]",
		},
		cli.StringFlag{
			Name:  "http-verb, h",
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
			Name:  "output-folder, o",
			Usage: "Set dir to write output files",
		},
		cli.StringFlag{
			Name:  "payload-key, p",
			Value: "data",
			Usage: "Specify the key to use for the payload when wrapping",
		},
		cli.StringFlag{
			Name:  "uuid, u",
			Usage: "Set namespaced uuid generation. uuid's will be injected with the 'id' key [full|wrap]",
		},
		cli.StringFlag{
			Name:  "wrapper, w",
			Usage: "Define wrapper a wrapper for the payload",
		},
		cli.BoolFlag{
			Name:  "recusrive, r",
			Usage: "Recursive read mode. also process sub-dirs",
		},
	}
}

func ProcessDocs(c *cli.Context) {
	// Create a buffered interface channel
	var dataChan = make(chan interface{}, c.Int("concurrency"))

	// Setup a new parser and pass the cli context and the dataChannel
	parser, err := NewParser(c, &dataChan)

	// Start a reader loop that reads data as quickly as possible and processes it in parallel

	// Parse all the input and dump the results
	contents, err := parser.Parse()

	// Start a writer loop that writes data as quickly as it comes available on the channel.
	// Should never exceed more then the concurrency configured. If write is slow this might
	// Increase memory consumption.

	// Generate output and write it to the api
	writer, err := NewWriter(c)
	writer.Write(contents)
}
