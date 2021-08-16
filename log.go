package gdcache

import (
	"fmt"
	"log"
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
	log.Default().Print(fmt.Sprintf(format, a))
}

func (d DefaultLogger) Error(format string, a ...interface{}) {
	log.Default().Print(fmt.Sprintf(format, a))
}

func (d DefaultLogger) Debug(format string, a ...interface{}) {
	log.Default().Print(fmt.Sprintf(format, a))
}

func (d DefaultLogger) Warn(format string, a ...interface{}) {
	log.Default().Print(fmt.Sprintf(format, a))
}
