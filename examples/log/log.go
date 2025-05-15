package log

import (
	"github.com/op/go-logging"
	"os"
)

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the millisecond.
const DefaultFormat = "%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.8s}%{color:reset} %{message}"

func NewLogger() *logging.Logger {
	return NewLoggerWithFormat(DefaultFormat)
}

func NewLoggerWithFormat(format string) *logging.Logger {
	formatter := logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0),
		logging.MustStringFormatter(format),
	)

	logging.SetBackend(formatter)

	return logging.MustGetLogger("goesl")
}
