package logger

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// TODO: this code should be moved to pkg, could be suggested as a PR for logrus
type logEntryCtxKeyType string

const logEntryCtxKey logEntryCtxKeyType = "log-entry"

// WithLogEntry attach a logrus entry to a context
func WithLogEntry(ctx context.Context, entry *log.Entry) context.Context {
	return context.WithValue(ctx, logEntryCtxKey, entry)
}

// GetLogEntry return logrus entry atteched to context
func GetLogEntry(ctx context.Context) *log.Entry {
	entry, ok := ctx.Value(logEntryCtxKey).(*log.Entry)
	if !ok {
		return log.NewEntry(log.StandardLogger())
	}
	return entry
}
