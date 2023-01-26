package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

func init() {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
}

func Info(v ...interface{}) {
	log.Info(v...)
}

func Warn(v ...interface{}) {
	log.Warn(v...)
}

func Error(v ...interface{}) {
	log.Error(v...)
}

func Debug(v ...interface{}) {
	log.Debug(v...)
}

func SetLevel(level string) {
	l, _ := logrus.ParseLevel(level)
	log.SetLevel(l)
}
