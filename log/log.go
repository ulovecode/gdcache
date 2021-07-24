package log

import (
	"fmt"
	defaultLog "log"
)

type Logger interface {
	Info(format string, a ...interface{})
	Error(format string, a ...interface{})
	Debug(format string, a ...interface{})
	Warn(format string, a ...interface{})
}

type DefaultLogger struct {
}

func (d DefaultLogger) Info(format string, a ...interface{}) {
	defaultLog.Default().Print(fmt.Sprintf(format, a))
}

func (d DefaultLogger) Error(format string, a ...interface{}) {
	defaultLog.Default().Print(fmt.Sprintf(format, a))
}

func (d DefaultLogger) Debug(format string, a ...interface{}) {
	defaultLog.Default().Print(fmt.Sprintf(format, a))
}

func (d DefaultLogger) Warn(format string, a ...interface{}) {
	defaultLog.Default().Print(fmt.Sprintf(format, a))
}
