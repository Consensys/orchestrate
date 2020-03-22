package testutils

import "context"

type Listener struct {
	Calls chan []string
}

func (l *Listener) Listen(ctx context.Context, mergedConf interface{}) error {
	l.Calls <- mergedConf.([]string)
	return nil
}
