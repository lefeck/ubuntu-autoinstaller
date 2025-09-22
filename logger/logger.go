package logger

import (
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// LogConfig holds logging configuration
type LogConfig struct {
	ShowCommandOutput bool
	CommandLogLevel   logrus.Level
}

var (
	Info   func(...interface{})
	Warn   func(...interface{})
	Error  func(...interface{})
	Infof  func(string, ...interface{})
	Warnf  func(string, ...interface{})
	Errorf func(string, ...interface{})
	Fatalf func(string, ...interface{})

	// Global log configuration
	Config = LogConfig{
		ShowCommandOutput: false,            // Set to false to hide command output
		CommandLogLevel:   logrus.InfoLevel, // Level for command execution logs
	}
)

func init() {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	Logger.SetLevel(logrus.InfoLevel)

	// Initialize function variables after Logger is created
	Info = Logger.Info
	Warn = Logger.Warn
	Error = Logger.Error
	Infof = Logger.Infof
	Warnf = Logger.Warnf
	Errorf = Logger.Errorf
	Fatalf = Logger.Fatalf
}

// SetCommandOutputEnabled enables or disables command output logging
func SetCommandOutputEnabled(enabled bool) {
	Config.ShowCommandOutput = enabled
}

// SetCommandLogLevel sets the log level for command execution
func SetCommandLogLevel(level logrus.Level) {
	Config.CommandLogLevel = level
}
