package util

import "log"

// Logger defines the logging interface used across the project.
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// StdLogger is a Logger backed by the standard library log package.
type StdLogger struct{}

// NewLogger returns a Logger that writes to stdout/stderr.
func NewLogger() Logger {
	return &StdLogger{}
}

// Infof logs an informational message.
func (l *StdLogger) Infof(format string, args ...interface{}) {
	log.Printf("INFO: "+format, args...)
}

// Errorf logs an error message.
func (l *StdLogger) Errorf(format string, args ...interface{}) {
	log.Printf("ERROR: "+format, args...)
}
