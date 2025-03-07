package utils

import (
	"fmt"

	"github.com/ghettovoice/gosip/log"
	"github.com/puzpuzpuz/xsync/v3"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type MyLogger struct {
	Logger *log.LogrusLogger
	level  log.Level
}

func (ml *MyLogger) Level() string {
	switch ml.level {
	case log.PanicLevel:
		return "Panic"
	case log.FatalLevel:
		return "Fatal"
	case log.ErrorLevel:
		return "Error"
	case log.WarnLevel:
		return "Warn"
	case log.InfoLevel:
		return "Info"
	case log.DebugLevel:
		return "Debug"
	case log.TraceLevel:
		return "Trace"
	}
	return "Unkown"
}

var (
	loggers *xsync.MapOf[string, *MyLogger]
)

func init() {
	loggers = xsync.NewMapOf[string, *MyLogger]()
}

func NewLogrusLogger(level log.Level, prefix string, fields log.Fields) log.Logger {
	if logger, found := loggers.Load(prefix); found {
		return logger.Logger.WithPrefix(prefix)
	}
	l := logrus.New()
	l.Level = logrus.ErrorLevel
	l.Formatter = &prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		ForceColors:     true,
		ForceFormatting: true,
	}
	l.SetReportCaller(true)
	logger := log.NewLogrusLogger(l, "main", fields)
	loggers.Store(prefix, &MyLogger{
		Logger: logger,
		level:  level,
	})
	logger.SetLevel(uint32((level)))
	return logger.WithPrefix(prefix)
}

func SetLogLevel(prefix string, level log.Level) error {
	if logger, found := loggers.Load(prefix); found {
		logger.level = level
		logger.Logger.SetLevel(uint32(level))
		return nil
	}
	return fmt.Errorf("logger [%v] not found", prefix)
}

func GetLoggers() *xsync.MapOf[string, *MyLogger] {
	return loggers
}
