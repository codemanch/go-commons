package logging

import (
	"testing"
)

var logger *Logger

var loggerTests = []struct {
	val string
}{
	{"12345"},
	{"testing"},
	{"\xff\xf0\x0f\xff"},
	{""},
	{"\""},
	{`\n`},
	{"\n"},
	{"日本語"},
	{"☺"},
	{"⌘"},
	{"\U0010ffff"},
}

// TestGetLogger --> Testing Logger object creation
func TestGetLogger(t *testing.T) {
	logger = GetLogger()
	if logger == nil {
		t.Errorf("Logger Object Test Fail!")
	}
}

func TestInfo(t *testing.T) {
	for _, tt := range loggerTests {
		logger.Info(tt.val)
	}
}

// TestGetLogger --> Testing InfoF
func TestInfoF(t *testing.T) {
	for _, tt := range loggerTests {
		logger.InfoF(tt.val)
	}
}

// TestGetLogger --> Testing Debug
func TestDebug(t *testing.T) {
	for _, tt := range loggerTests {
		logger.Debug(tt.val)
	}
}

// TestGetLogger --> Testing DebugF
func TestDebugF(t *testing.T) {
	for _, tt := range loggerTests {
		logger.DebugF(tt.val)
	}
}

// TestGetLogger --> Testing Trace
func TestTrace(t *testing.T) {
	for _, tt := range loggerTests {
		logger.Trace(tt.val)
	}
}

// TestGetLogger --> Testing TraceF
func TestTraceF(t *testing.T) {
	for _, tt := range loggerTests {
		logger.TraceF(tt.val)
	}
}

// TestGetLogger --> Testing Warn
func TestWarn(t *testing.T) {
	for _, tt := range loggerTests {
		logger.Warn(tt.val)
	}
}

// TestGetLogger --> Testing WarnF
func TestWarnF(t *testing.T) {
	for _, tt := range loggerTests {
		logger.WarnF(tt.val)
	}
}

// TestGetLogger --> Testing Error
func TestError(t *testing.T) {
	for _, tt := range loggerTests {
		logger.Error(tt.val)
	}
}

// TestGetLogger --> Testing ErrorF
func TestErrorF(t *testing.T) {
	for _, tt := range loggerTests {
		logger.ErrorF(tt.val)
	}
}
