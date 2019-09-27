package jaeger

import (
	log "github.com/sirupsen/logrus"
)

// logger is an adpatator from logrus to jaeger.Logger
type logger struct {
	entry *log.Entry
}

func (l logger) Error(msg string) {
	l.entry.Error(msg)
}

func (l logger) Infof(msg string, args ...interface{}) {
	// Jaegger logger Infof is called each time a span is reported
	// so we prefer to debug to avoid logging to much
	l.entry.Debugf(msg, args...)
}
