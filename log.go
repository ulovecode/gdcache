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

// Info Info category log
func (d DefaultLogger) Info(format string, a ...interface{}) {
	fmt.Print(fmt.Sprintf(format, a...))
}

// Error Error category log
func (d DefaultLogger) Error(format string, a ...interface{}) {
	fmt.Print(fmt.Sprintf(format, a...))
}

// Debug Debug category log
func (d DefaultLogger) Debug(format string, a ...interface{}) {
	fmt.Print(fmt.Sprintf(format, a...))
}

// Warn Warn category log
func (d DefaultLogger) Warn(format string, a ...interface{}) {
	fmt.Print(fmt.Sprintf(format, a...))
}
