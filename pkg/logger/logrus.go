package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Init() *logrus.Logger {
	// Create a new logger instance
	logger := logrus.New()

	// Set the output for the logger to stdout
	logger.SetOutput(os.Stdout)

	// Set the log level. You can change this to the desired level, e.g., logrus.DebugLevel for more verbose logs.
	logger.SetLevel(logrus.InfoLevel)

	// Optionally, you can set a formatter to customize the log output.
	// For example, to display the timestamp and the log level in the log messages:
	// logger.SetFormatter(&logrus.TextFormatter{
	// 	FullTimestamp: true,
	// })

	return logger
}
