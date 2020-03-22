package grpclogrus

import (
	"github.com/sirupsen/logrus"
)

// Entry is a wrapper around logrus Entry that implements
// grpclog.LoggerV2 interface
type Entry struct {
	*logrus.Entry
}

// V reports whether verbosity level l is at least the requested verbose level.
func (e *Entry) V(l int) bool {
	// We can always answer true, as verbosity level is already checked by logrus
	return true
}
