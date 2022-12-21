package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Formatter = &logrus.TextFormatter{
		TimestampFormat: "15:04:05",
		FullTimestamp:   true,
		ForceColors:     true,
		DisableColors:   false,
	}

	log.Out = os.Stdout
}

// SetLogLevel sets the log level for the logger
func SetLogLevel(level string) {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		log.Errorf("Invalid log level: %s, using default level: info", level)
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)

}

func WriteToFileAndConsole(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	mw := io.MultiWriter(os.Stdout, file)
	log.Out = mw

	return nil
}

// WriteToFile writes the log output to a file
func WriteToFile(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	log.Out = file
	return nil
}

// Debug logs a message with debug level
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Info logs a message with info level
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn logs a message with warning level
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Error logs a message with error level
func Error(args ...interface{}) {
	log.Error(args...)
}

// Fatal logs a message with fatal level and then calls os.Exit(1)
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Panic logs a message with panic level and then calls panic()
func Panic(args ...interface{}) {
	log.Panic(args...)
}
