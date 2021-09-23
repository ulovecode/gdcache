package gdcache

import (
	"fmt"
)

// Logger Log component
type Logger interface {
	// Info Info category log
	Info(format string, a ...interface{})
	// Error Error category log
	Error(format string, a ...interface{})
	// Debug Debug category log
	Debug(format string, a ...interface{})
	// Warn Warn category log
	Warn(format string, a ...interface{})
}

// DefaultLogger Default log component
type DefaultLogger struct {
}

func (d DefaultLogger) Info(format string, a ...interface{}) {
	fmt.Print(fmt.Sprintf(format, a))
}

func (d DefaultLogger) Error(format string, a ...interface{}) {
	fmt.Print(fmt.Sprintf(format, a))
}

func (d DefaultLogger) Debug(format string, a ...interface{}) {
	fmt.Print(fmt.Sprintf(format, a))
}

func (d DefaultLogger) Warn(format string, a ...interface{}) {
	fmt.Print(fmt.Sprintf(format, a))
}
