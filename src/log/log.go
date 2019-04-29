package log

import (
	"github.com/sirupsen/logrus"
)

var entry *logrus.Entry

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   false,
		TimestampFormat: "2006-01-02 15:04:05.00",
	})
	entry = logrus.WithFields(logrus.Fields{"Origin": "System Log"})
}

func Info(v ...interface{}) {
	entry.Info(v...)
}

func Infof(format string, v ...interface{}) {
	entry.Infof(format, v...)
}

func Error(v ...interface{}) {
	entry.Error(v)
}

func Errorf(format string, v ...interface{}) {
	entry.Errorf(format, v...)
}

func Debug(v ...interface{}) {
	entry.Debug(v)
}

func Debugf(format string, v ...interface{}) {
	entry.Debugf(format, v...)
}

func Panic(v ...interface{}) {
	entry.Panic(v...)
}