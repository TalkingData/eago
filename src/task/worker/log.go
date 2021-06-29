package worker

import "fmt"

type Logger interface {
	Info(format string, a ...interface{})
	Error(format string, a ...interface{})
}

type logger struct {
	LogCh chan *string
}

func (l *logger) Info(format string, a ...interface{}) {
	msg := "[INFO] " + fmt.Sprintf(format, a...)
	fmt.Println(msg)
	go l.putMsg(msg)
}

func (l *logger) Error(format string, a ...interface{}) {
	msg := "[ERROR] " + fmt.Sprintf(format, a...)
	fmt.Println(msg)
	go l.putMsg(msg)
}

func (l *logger) putMsg(msg string) {
	l.LogCh <- &msg
}

func NewLogger(bufferSize uint) *logger {
	l := &logger{}
	l.LogCh = make(chan *string, bufferSize)
	return l
}
