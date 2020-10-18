package logging

import "time"

// LogMessage struct.
type LogMessage struct {
	Time   time.Time `json:"timestamp"`
	FnName string    `json:"function,omitempty"`
	Line   int       `json:"line,omitempty"`
	Msg    string    `json:"msg"`
	Sev    Severity  `json:"sev"`
}
