package util

import (
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/npolar/ghostdoc/context"
)

// ConfigureLog4go configures log4go
func ConfigureLog4go(c context.GhostContext) {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	if logfile := c.GlobalString("logfile"); logfile != "" {
		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}

		defer f.Close()
		log.SetOutput(f)
	}
	if c.GlobalBool("quiet") {
		log.SetOutput(ioutil.Discard)
	}

	if c.GlobalBool("no-verbose") {
		log.SetLevel(log.InfoLevel)
	}
	log.Debug("Logging configured")
}
