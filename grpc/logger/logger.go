package grpclogger

import (
	log "github.com/sirupsen/logrus"
)

// LogEntry is a wrapper around logrus Entry that implements
// grpclog.LoggerV2 interface
type LogEntry struct {
	*log.Entry
}

// V reports whether verbosity level l is at least the requested verbose level.
func (e *LogEntry) V(l int) bool {
	// We can always answer true, as verbosity level is already checked by logrus
	return true
}
