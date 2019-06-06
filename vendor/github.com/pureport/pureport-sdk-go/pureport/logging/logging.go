package logging

import (
	"io"
	"os"
	"strings"

	"github.com/op/go-logging"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig contains setting for configuring the logger
type LogConfig struct {
	Level       string `desc:"One of the available log log-levels- (info, debug, error, warn)"`
	Location    string `desc:"stdout, stderr, or the path to a file"`
	MaxSize     int    `desc:"Maximum size (in MB) of the log file before it gets rotated, default=10"`
	MaxBackups  int    `desc:"Maximum number of old log files to retain, default=2"`
	MaxAge      int    `desc:"Maximum number of days to retain old log files, default=1"`
	EnableColor bool   `desc:"Enable color log output"`
}

var log = logging.MustGetLogger("main_logger")

const (
	// LogLevelEnvStr - environment variable for the setting the log level
	LogLevelEnvStr string = "PUREPORT_LOG_LEVEL"

	// LogFileEnvStr - environment variable for the location of the log file
	LogFileEnvStr string = "PUREPORT_LOG_FILE"

	// LogDisableColorStr - environment variable to disable color log output
	LogDisableColorStr string = "PUREPORT_LOG_NOCOLOR"
)

// NewLogConfig generate a new default logging configuration
func NewLogConfig() *LogConfig {
	return &LogConfig{
		Level:       "info",
		Location:    "stdout",
		MaxSize:     10, // in MB
		MaxBackups:  2,
		MaxAge:      1, // in days
		EnableColor: true,
	}
}

// SetupLogger configures the system logger based on the specified
// log configuration. Environment variables can override these values.
func SetupLogger(logConfig *LogConfig) {

	if color := os.Getenv(LogDisableColorStr); len(color) != 0 {
		logConfig.EnableColor = false
	}

	if location := os.Getenv(LogFileEnvStr); len(location) != 0 {
		logConfig.Location = location
	}

	if level := os.Getenv(LogLevelEnvStr); len(level) != 0 {
		logConfig.Level = level
	}

	// Initialize the logger
	format := logging.MustStringFormatter(
		`%{color}%{time:2006-01-02T15:04:05.000000Z07:00} %{shortfile} %{shortfunc} ▶ %{level:.4s} %{id:03d}%{color:reset} %{message}`,
	)

	// Disable color is enabled
	if !logConfig.EnableColor {
		format = logging.MustStringFormatter(
			`%{time:2006-01-02T15:04:05.000000Z07:00} %{shortfile} %{shortfunc} ▶ %{level:.4s} %{id:03d} %{message}`,
		)
	}

	//LOG LOCATION
	var logOut io.Writer = os.Stdout

	switch strings.ToLower(logConfig.Location) {

	case "stdout", "":
		logOut = os.Stdout

	case "stderr":
		logOut = os.Stderr

	default:
		// Setup Log Rotation
		logOut = &lumberjack.Logger{
			Filename:   logConfig.Location,
			MaxSize:    logConfig.MaxSize, // in MB
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge, // in days
			Compress:   true,
		}
	}

	backend := logging.NewLogBackend(logOut, "", 0)

	level, err := logging.LogLevel(logConfig.Level)
	if err != nil {
		level = logging.INFO
	}

	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
	logging.SetLevel(level, "main_logger")
}
