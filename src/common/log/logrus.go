package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

var logger *logrus.Logger

type Fields logrus.Fields

func InitLog(logPath string, srvName string) error {
	logger = logrus.New()
	//logger.SetReportCaller(true)

	logPathName := filepath.Join(logPath, srvName+".log")

	fw := &lumberjack.Logger{
		Filename:  logPathName,
		LocalTime: true,
		MaxSize:   20,
		MaxAge:    31,
	}
	logger.SetOutput(io.MultiWriter(os.Stdout, fw))

	// 保留调试用
	//formatter := &logrus.TextFormatter{}

	formatter := &logrus.JSONFormatter{}
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	logger.SetFormatter(formatter)

	logger.SetLevel(logrus.InfoLevel)
	return nil
}

// Debug
func Debug(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Debug(args)
	}
}

// 带有field的Debug
func DebugWithFields(f Fields, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Debug(args...)
	}
}

// Info
func Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Info(args...)
	}
}

// 带有field的Info
func InfoWithFields(f Fields, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Info(args...)
	}
}

// Warn
func Warn(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Warn(args...)
	}
}

// 带有Field的Warn
func WarnWithFields(f Fields, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Warn(args...)
	}
}

// Error
func Error(args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Error(args...)
	}
}

// 带有Fields的Error
func ErrorWithFields(f Fields, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Error(args...)
	}
}

// Fatal
func Fatal(args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(args...)
	}
}

// 带有Field的Fatal
func FatalWithFields(f Fields, args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(args...)
	}
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = -1
	}
	return fmt.Sprintf("%s:%d", file, line)
}
