package worker

import (
	"fmt"
	"sync"
)

type ResultLog interface {
	Debug(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warn(format string, a ...interface{})
	Error(format string, a ...interface{})
}

type resultLog struct {
	printLog bool

	wg    *sync.WaitGroup
	logCh chan *string
}

func (resLog *resultLog) Debug(format string, a ...interface{}) {
	resLog.wg.Add(1)
	content := "[DEBUG] " + fmt.Sprintf(format, a...)
	if resLog.printLog {
		fmt.Println(content)
	}
	resLog.logCh <- &content
}

func (resLog *resultLog) Info(format string, a ...interface{}) {
	resLog.wg.Add(1)
	content := "[INFO] " + fmt.Sprintf(format, a...)
	if resLog.printLog {
		fmt.Println(content)
	}
	resLog.logCh <- &content
}

func (resLog *resultLog) Warn(format string, a ...interface{}) {
	resLog.wg.Add(1)
	content := "[WARNING] " + fmt.Sprintf(format, a...)
	if resLog.printLog {
		fmt.Println(content)
	}
	resLog.logCh <- &content
}

// Error
func (resLog *resultLog) Error(format string, a ...interface{}) {
	resLog.wg.Add(1)
	content := "[ERROR] " + fmt.Sprintf(format, a...)
	if resLog.printLog {
		fmt.Println(content)
	}
	resLog.logCh <- &content
}

// newResultLog 创建一个ResultLog
func newResultLog(bufferSize uint, printLog bool) *resultLog {
	return &resultLog{
		printLog: printLog,

		wg:    &sync.WaitGroup{},
		logCh: make(chan *string, bufferSize),
	}
}
