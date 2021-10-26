package worker

import (
	"fmt"
	"sync"
)

type Logger interface {
	Debug(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warn(format string, a ...interface{})
	Error(format string, a ...interface{})
}

type logger struct {
	Wg    *sync.WaitGroup
	LogCh chan *string
}

// Debug
func (l *logger) Debug(format string, a ...interface{}) {
	l.Wg.Add(1)
	msg := "[DEBUG] " + fmt.Sprintf(format, a...)
	fmt.Println(msg)
	l.LogCh <- &msg
}

// Info
func (l *logger) Info(format string, a ...interface{}) {
	l.Wg.Add(1)
	msg := "[INFO] " + fmt.Sprintf(format, a...)
	fmt.Println(msg)
	l.LogCh <- &msg
}

// Warn
func (l *logger) Warn(format string, a ...interface{}) {
	l.Wg.Add(1)
	msg := "[WARNING] " + fmt.Sprintf(format, a...)
	fmt.Println(msg)
	l.LogCh <- &msg
}

// Error
func (l *logger) Error(format string, a ...interface{}) {
	l.Wg.Add(1)
	msg := "[ERROR] " + fmt.Sprintf(format, a...)
	fmt.Println(msg)
	l.LogCh <- &msg
}

// newLogger 创建一个Logger
func newLogger(bufferSize uint) *logger {
	l := &logger{}
	l.Wg = &sync.WaitGroup{}
	l.LogCh = make(chan *string, bufferSize)
	return l
}
