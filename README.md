#ghostdoc - Flexible file parser / REST client

##USAGE:
   ghostdoc [global options] command [command options] [arguments...]

##VERSION:
   0.0.1

##COMMANDS:
   * csv: Parse delimiter separated value files (csv, tsv, etc...)
   * json: Parse json files
   * text, txt: Parse text data
   * help, h:	Shows a list of commands or help for one command

##GLOBAL OPTIONS:
   * --address, -a 		Set url to write to
   * --concurrency, -c "2"	Specify the number of concurrent operations
   * --exclude, -e 		Specify keys (before mapping) to exclude in the output
   * --filename, -f 		Set filename to use in name-pattern when piping data via stdin
   * --include, -i 		Specify keys (before mapping) to include in the output
   * --js, -j 			Run javascript map functions on the data
   * --http-verb "POST"		Set the http verb to use [POST|PUT]
   * --key-map, -k 		Sets mapping file to use to rename headers/keys. JSON Format {"oldkey": "newkey"}
   * --merge, -m 			Specify additional JSON data to inject into the output.
   * --name-pattern, -n 		Set pattern file to extract filename info and inject it into the result
   * --output, -o 		Set dir output dir. Files will get uuid as name
   * --payload-key, -p "data"	Specify the key to use for the payload when wrapping
   * --quiet, -q			Turn off logging to stdout
   * --uuid, -u			Injects a namesaced uuid with the 'id' key
   * --uuid-keys, --uk 		Injects a namesaced uuid with the 'id' key based on a set of keys
   * --wrapper, -w 		Define JSON wrapper a wrapper for the payload
   * --recursive, -r		Recursive read mode. also process sub-dirs
   * --help, -h			show help
   * --version, -v		print the version

## javascript api
javascript mapping files should expose one global object named "functions".

```
functions = {
  "fn1": function (doc) {
    doc.js = "JavaScript";
    return doc;
  },
  "fn2": function (doc) {
    delete doc.js;
    return doc;
  }
};
```
