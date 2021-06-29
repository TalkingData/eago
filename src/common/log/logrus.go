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

var (
	logger    *logrus.Logger
	skipFiles int = 2
)

type Fields logrus.Fields

// InitLog 初始化log
func InitLog(logPath, srvName, timestampFormatter string, lvl logrus.Level) error {
	logger = logrus.New()

	logPathName := filepath.Join(logPath, srvName+".log")

	fw := &lumberjack.Logger{
		Filename:  logPathName,
		LocalTime: true,
		MaxSize:   20,
		MaxAge:    31,
	}
	logger.SetOutput(io.MultiWriter(os.Stdout, fw))

	formatter := &logrus.JSONFormatter{}
	formatter.TimestampFormat = timestampFormatter
	logger.SetFormatter(formatter)

	logger.SetLevel(lvl)
	return nil
}

func SetSkipFiles(skip int) {
	skipFiles = skip
}

// Debug
func Debug(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Debug(args)
	}
}

// DebugWithFields
func DebugWithFields(f Fields, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Debug(args...)
	}
}

// Info
func Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Info(args...)
	}
}

// InfoWithFields
func InfoWithFields(f Fields, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Info(args...)
	}
}

// Warn
func Warn(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Warn(args...)
	}
}

// WarnWithFields
func WarnWithFields(f Fields, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Warn(args...)
	}
}

// Error
func Error(args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Error(args...)
	}
}

// ErrorWithFields
func ErrorWithFields(f Fields, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Error(args...)
	}
}

// Fatal
func Fatal(args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Fatal(args...)
	}
}

// FatalWithFields
func FatalWithFields(f Fields, args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(skipFiles)
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
