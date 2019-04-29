package log

import (
	"github.com/sirupsen/logrus"
)

func init() {
	//logrus.SetReportCaller(false)
	//logrus.SetOutput(os.Stdout)
	//logrus.SetLevel(logrus.DebugLevel)
}

func Info(v ...interface{}) {
	logrus.Info(v...)
}

func Infof(format string, v ...interface{}) {
	logrus.Infof(format, v...)
}

func Error(v ...interface{}) {
	logrus.Error(v)
}

func Errorf(format string, v ...interface{}) {
	logrus.Errorf(format, v...)
}

func Debug(v ...interface{}) {
	logrus.Debug(v)
}

func Debugf(format string, v ...interface{}) {
	logrus.Debugf(format, v...)
}

func Panic(v ...interface{}) {
	logrus.Panic(v...)
}
