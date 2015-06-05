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
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	if logfile := c.GlobalString("log-file"); logfile != "" {
		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
		log.SetOutput(f)
	}
	switch c.GlobalString("log-level") {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "OFF":
		log.SetOutput(ioutil.Discard)
	}

	log.Debug("Logging configured")
}
