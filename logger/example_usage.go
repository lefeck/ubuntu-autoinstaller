package logger

import (
	"github.com/sirupsen/logrus"
)

// Example usage of the new logging configuration
func ExampleUsage() {
	// Disable command output logging (useful for production)
	SetCommandOutputEnabled(false)

	// Or enable it but set to debug level (useful for development)
	SetCommandOutputEnabled(true)
	SetCommandLogLevel(logrus.DebugLevel)

	// Or set to warn level to only show warnings and errors
	SetCommandLogLevel(logrus.WarnLevel)

	// Or set to info level (default) to show all command execution
	SetCommandLogLevel(logrus.InfoLevel)
}

// Example of how to configure logging in main.go
func ConfigureLogging() {
	// For production: hide command output, only show errors
	SetCommandOutputEnabled(false)
	Logger.SetLevel(logrus.WarnLevel)

	// For development: show all command output
	SetCommandOutputEnabled(true)
	SetCommandLogLevel(logrus.InfoLevel)
	Logger.SetLevel(logrus.DebugLevel)
}
