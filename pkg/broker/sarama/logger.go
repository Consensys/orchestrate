package sarama

import (
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/log"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type logger struct {
	*log.Logger
	level logrus.Level
}

var _ sarama.StdLogger = logger{}

func newLogger(l *log.Logger) *logger {
	return &logger{l, logrus.InfoLevel}
}

func (l logger) SetLevel(lvl logrus.Level) *logger {
	l.level = lvl
	return &l
}

func (l logger) Print(v ...interface{}) {
	switch l.level {
	case logrus.InfoLevel:
		l.Info(v...)
	case logrus.DebugLevel:
		l.Debug(v...)
	case logrus.TraceLevel:
		l.Trace(v...)
	}
}

func (l logger) Printf(format string, v ...interface{}) {
	switch l.level {
	case logrus.InfoLevel:
		l.Info(fmt.Sprintf(format, v...))
	case logrus.DebugLevel:
		l.Debug(fmt.Sprintf(format, v...))
	case logrus.TraceLevel:
		l.Trace(fmt.Sprintf(format, v...))
	}
}

func (l logger) Println(v ...interface{}) {
	switch l.level {
	case logrus.InfoLevel:
		l.Info(v...)
	case logrus.DebugLevel:
		l.Debug(v...)
	case logrus.TraceLevel:
		l.Trace(v...)
	}
}
