package logging

import (
	"bytes"
	"fmt"
	"github.com/appmanch/go-commons/textutils"
	"sync"
	"time"
)

var logMsgPool = &sync.Pool{
	New: func() interface{} {
		return &LogMessage{
			Buf: &bytes.Buffer{},
		}
	},
}

// LogMessage struct.
type LogMessage struct {
	Time     time.Time     `json:"timestamp"`
	FnName   string        `json:"function,omitempty"`
	Line     int           `json:"line,omitempty"`
	Buf      *bytes.Buffer `json:"msg"`
	Sev      Severity      `json:"sev"`
	SevBytes []byte
}

func getLogMessageF(sev Severity, f string, v ...interface{}) *LogMessage {
	msg := logMsgPool.Get().(*LogMessage)
	msg.Sev = sev
	msg.Time = time.Now()
	msg.FnName = textutils.EmptyStr
	msg.Line = 0
	fmt.Fprintf(msg.Buf, f, v...)
	return msg
}

func getLogMessage(sev Severity, v ...interface{}) *LogMessage {
	msg := logMsgPool.Get().(*LogMessage)
	msg.Sev = sev
	msg.Time = time.Now()
	msg.FnName = textutils.EmptyStr
	msg.Line = 0
	fmt.Fprint(msg.Buf, v...)
	return msg

}

func putLogMessage(logMsg *LogMessage) {
	logMsg.Buf.Reset()
	logMsgPool.Put(logMsg)
}
