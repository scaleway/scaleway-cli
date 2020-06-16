package core

import (
	"fmt"
	"io"

	"github.com/scaleway/scaleway-sdk-go/logger"
)

type Logger struct {
	writer io.Writer
	level  logger.LogLevel
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.ShouldLog(logger.LogLevelDebug) {
		_, _ = fmt.Fprintf(l.writer, format, args...)
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.ShouldLog(logger.LogLevelInfo) {
		_, _ = fmt.Fprintf(l.writer, format, args...)
	}
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	if l.ShouldLog(logger.LogLevelWarning) {
		_, _ = fmt.Fprintf(l.writer, format, args...)
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.ShouldLog(logger.LogLevelError) {
		_, _ = fmt.Fprintf(l.writer, format, args...)
	}
}

func (l *Logger) Debug(args ...interface{}) {
	if l.ShouldLog(logger.LogLevelDebug) {
		_, _ = fmt.Fprintln(l.writer, args...)
	}
}

func (l *Logger) Info(args ...interface{}) {
	if l.ShouldLog(logger.LogLevelInfo) {
		_, _ = fmt.Fprintln(l.writer, args...)
	}
}

func (l *Logger) Warning(args ...interface{}) {
	if l.ShouldLog(logger.LogLevelWarning) {
		_, _ = fmt.Fprintln(l.writer, args...)
	}
}

func (l *Logger) Error(args ...interface{}) {
	if l.ShouldLog(logger.LogLevelError) {
		_, _ = fmt.Fprintln(l.writer, args...)
	}
}

func (l *Logger) ShouldLog(level logger.LogLevel) bool {
	return l.level <= level
}
