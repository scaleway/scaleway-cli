package clilogger

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/moul/http2curl"
	"github.com/scaleway/scaleway-cli/pkg/api"
)

type cliLogger struct {
	*logrus.Logger
	s *api.ScalewayAPI
}

func (l *cliLogger) LogHTTP(req *http.Request) {
	curl, err := http2curl.GetCurlCommand(req)
	if err != nil {
		l.Fatalf("Failed to convert to curl request: %q", err)
	}

	if os.Getenv("SCW_SENSITIVE") != "1" {
		l.Debug(l.s.HideAPICredentials(curl.String()))
	} else {
		l.Debug(curl.String())
	}
}

func NewCliLogger(s *api.ScalewayAPI) api.Logger {
	return &cliLogger{
		Logger: logrus.StandardLogger(),
		s:      s,
	}
}

func SetupLogger(s *api.ScalewayAPI) {
	s.Logger = NewCliLogger(s)
}
