package alert

import (
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

const server = "http://alert"

func send(msg, path string) error {
	u, err := url.Parse(server)
	if err != nil {
		return err
	}

	u.Path = path
	query := u.Query()
	query.Set("msg", msg)
	u.RawQuery = query.Encode()

	_, err = http.Post(u.String(), "text/plain", nil)
	return err
}

// Alert sends an alert to the log and audio alerting system
func Alert(msg string) {
	if err := send(msg, "/alert"); err != nil {
		log.Warnf("sending alert: %s", err)
	}
}

// Log sends an alert to the log only
func Log(msg string) {
	if err := send(msg, "/log"); err != nil {
		log.Warnf("sending log: %s", err)
	}
}
