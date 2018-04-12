package logger

import (
	"log"
	"os"

	"github.com/abhishekkr/gol/golenv"
	"github.com/sirupsen/logrus"
)

var ogilogger *logrus.Logger

func SetupLogger() {
	level, err := logrus.ParseLevel(golenv.OverrideIfEnv("LOG_LEVEL", "info"))
	if err != nil {
		log.Fatal(err)
	}

	ogilogger = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}
}

func Debug(args ...interface{}) {
	ogilogger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	ogilogger.Debugf(format, args...)
}

func Debugln(args ...interface{}) {
	ogilogger.Debugln(args...)
}

func Error(args ...interface{}) {
	ogilogger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	ogilogger.Errorf(format, args...)
}

func Errorln(args ...interface{}) {
	ogilogger.Errorln(args...)
}

func Fatal(args ...interface{}) {
	ogilogger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	ogilogger.Fatalf(format, args...)
}

func Fatalln(args ...interface{}) {
	ogilogger.Fatalln(args...)
}

func Info(args ...interface{}) {
	ogilogger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	ogilogger.Infof(format, args...)
}

func Infoln(args ...interface{}) {
	ogilogger.Infoln(args...)
}

func Warn(args ...interface{}) {
	ogilogger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	ogilogger.Warnf(format, args...)
}

func Warnln(args ...interface{}) {
	ogilogger.Warnln(args...)
}
