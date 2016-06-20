// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package api

import (
	"log"
	"net/http"
	"os"
)

// Logger handles logging concerns for the Scaleway API SDK
type Logger interface {
	LogHTTP(*http.Request)
	Log(...interface{})
}

// NewDefaultLogger returns a logger which is configured for stdout
func NewDefaultLogger() Logger {
	return &defaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

type defaultLogger struct {
	logger *log.Logger
}

func (l defaultLogger) LogHTTP(r *http.Request) {
	l.logger.Printf("%s %s", r.Method, r.URL.RawPath)
}

func (l defaultLogger) Log(args ...interface{}) {
	l.logger.Println(args...)
}
