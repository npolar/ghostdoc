package util

import (
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"os/user"

	log "github.com/Sirupsen/logrus"
	"github.com/npolar/ghostdoc/context"
)

type errorHook struct {
	to       string
	from     string
	message  string
	hasError bool
}

func newErrorHook(to string) *errorHook {
	user, _ := user.Current()
	host, _ := os.Hostname()
	from := user.Username + "@" + host

	header := make(map[string]string)
	header["From"] = from
	header["To"] = to
	header["Subject"] = "ghostdoc error!"
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	return &errorHook{to: to, message: message, from: from, hasError: false}
}

// Levels impl
func (eh *errorHook) Levels() []log.Level {
	return []log.Level{log.ErrorLevel}
}

// Fire impl
func (eh *errorHook) Fire(e *log.Entry) error {
	buf, _ := e.Reader()
	eh.message += "\n" + buf.String()
	return nil
}

var eh *errorHook

// SendErrorMail flush message to mail
func SendErrorMail() {
	if eh != nil && eh.hasError {
		err := smtp.SendMail("localhost:25", nil, eh.from, []string{eh.to}, []byte(eh.message))
		if err != nil {
			log.Warn(err)
		}
	}
}

// ConfigureLogger configures logging
func ConfigureLogger(c context.GhostContext) {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})

	if to := c.GlobalString("log-mail"); to != "" {
		eh = newErrorHook(to)
		log.AddHook(eh)
	}

	if logfile := c.GlobalString("log-file"); logfile != "" {
		f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
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
