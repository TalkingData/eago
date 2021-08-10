package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	LOG_DATA_SEPARATOR = "; "
)

var (
	logger    *logrus.Logger
	skipFiles int = 2
	fw        *lumberjack.Logger
)

type Fields logrus.Fields

// InitLog 初始化log
func InitLog(path, srvName, level string) error {
	logger = logrus.New()

	logPathName := filepath.Join(path, srvName+".log")

	fw = &lumberjack.Logger{
		Filename:  logPathName,
		LocalTime: true,
		MaxSize:   10,
		MaxAge:    3,
	}
	logger.SetOutput(io.MultiWriter(os.Stdout, fw))

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logLvl, err := logrus.ParseLevel(level)
	if err != nil {
		panic(err)
	}
	logger.SetLevel(logLvl)
	return nil
}

func Close() {
	if fw == nil {
		return
	}
	_ = fw.Close()
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
		entry := logger.WithFields(f.RusFields())
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
		entry := logger.WithFields(f.RusFields())
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
		entry := logger.WithFields(f.RusFields())
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
		entry := logger.WithFields(f.RusFields())
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
		entry := logger.WithFields(f.RusFields())
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Fatal(args...)
	}
}

// Panic
func Panic(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	entry.Data["file"] = fileInfo(skipFiles)
	entry.Panic(args...)
}

// PanicWithFields
func PanicWithFields(f Fields, args ...interface{}) {
	entry := logger.WithFields(f.RusFields())
	entry.Data["file"] = fileInfo(skipFiles)
	entry.Panic(args...)
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = -1
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// String 将logger.Fields转换为字符串形式
func (f Fields) String() string {
	data := []string{}
	for k, v := range f {
		data = append(data, fmt.Sprintf("%s=%v", k, v))
	}

	return strings.Join(data, LOG_DATA_SEPARATOR)
}

// 将logger.Fields转换为logrus.Fields
func (f Fields) RusFields() logrus.Fields {
	res := logrus.Fields{}

	if f["error"] != nil {
		res["error"] = f["error"]
		delete(f, "error")
	}

	if f["trace_id"] != nil {
		res["trace_id"] = f["trace_id"]
		delete(f, "trace_id")
	}

	res["data"] = f.String()

	return res
}
