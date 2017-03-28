package log

import (
	"github.com/Sirupsen/logrus"
	"os"
)

var logger = logrus.New()

func SetOutput() {
	logger.Out = os.Stdout
}

func SetLevel() {
	logger.Level = logrus.DebugLevel
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}
