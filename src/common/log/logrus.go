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
	MAX_LOG_SIZE       = 3
	MAX_LOG_AGE        = 5
	TIMESTAMP_FORMAT   = "2006-01-02 15:04:05"
	LOG_DATA_SEPARATOR = "; "
	SRC_BASE_PATH      = "/td-eago/src/"
)

var (
	logger    *logrus.Logger
	skipFiles = 2
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
		MaxSize:   MAX_LOG_SIZE,
		MaxAge:    MAX_LOG_AGE,
	}
	logger.SetOutput(io.MultiWriter(os.Stdout, fw))

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: TIMESTAMP_FORMAT,
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

func Debug(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Debug(args)
	}
}

func DebugWithFields(f Fields, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(f.RusFields())
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Debug(args...)
	}
}

func Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Info(args...)
	}
}

func InfoWithFields(f Fields, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(f.RusFields())
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Info(args...)
	}
}

func Warn(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Warn(args...)
	}
}

func WarnWithFields(f Fields, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(f.RusFields())
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Warn(args...)
	}
}

func Error(args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Error(args...)
	}
}

func ErrorWithFields(f Fields, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(f.RusFields())
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Error(args...)
	}
}

func Fatal(args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Fatal(args...)
	}
}

func FatalWithFields(f Fields, args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(f.RusFields())
		entry.Data["file"] = fileInfo(skipFiles)
		entry.Fatal(args...)
	}
}

func Panic(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	entry.Data["file"] = fileInfo(skipFiles)
	entry.Panic(args...)
}

func PanicWithFields(f Fields, args ...interface{}) {
	entry := logger.WithFields(f.RusFields())
	entry.Data["file"] = fileInfo(skipFiles)
	entry.Panic(args...)
}

// fileInfo
func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "<Unknown>:-1"
	}

	file = file[strings.Index(file, SRC_BASE_PATH)+len(SRC_BASE_PATH):]
	return fmt.Sprintf("./%s:%d", file, line)
}

// String 将logger.Fields转换为字符串形式
func (f Fields) String() string {
	var data []string
	for k, v := range f {
		data = append(data, fmt.Sprintf("%s=%v", k, v))
	}

	return strings.Join(data, LOG_DATA_SEPARATOR)
}

// RusFields 将logger.Fields转换为logrus.Fields
func (f Fields) RusFields() logrus.Fields {
	res := logrus.Fields{}

	// 从data中提取code字段，code字段不会出现在最终输入的data中
	if f["code"] != nil {
		res["code"] = f["code"]
		delete(f, "code")
	}

	// 从data中提取error字段，error字段不会出现在最终输入的data中
	if f["error"] != nil {
		switch errType := f["error"].(type) {
		// 如果是error类型，设置最终error字段为err.Error()
		case error:
			res["error"] = errType.Error()
		// 如果是其他类型，设置最终error字段为error（不变）
		default:
			res["error"] = f["error"]
		}
		delete(f, "error")
	}

	// 从data中提取trace_id字段，trace_id字段不会出现在最终输入的data中
	if f["trace_id"] != nil {
		res["trace_id"] = f["trace_id"]
		delete(f, "trace_id")
	}

	// 将data转为字符串形式
	res["data"] = f.String()

	return res
}
