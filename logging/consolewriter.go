package logging

import (
	"io"
	"os"
)

// ConsoleWriter struct
type ConsoleWriter struct {
	errorWriter, warnWriter, infoWriter, debugWriter, traceWriter io.Writer
}

// InitConfig Consolewriter
func (cw *ConsoleWriter) InitConfig(w *WriterConfig) {
	if w.Console.WriteErrToStdOut {
		cw.errorWriter = os.Stdout
	} else {
		cw.errorWriter = os.Stderr
	}
	if w.Console.WriteWarnToStdOut {
		cw.warnWriter = os.Stdout
	} else {
		cw.warnWriter = os.Stderr
	}

	cw.infoWriter = os.Stdout
	cw.debugWriter = os.Stdout
	cw.traceWriter = os.Stdout

}

// DoLog consolewriter
func (cw *ConsoleWriter) DoLog(logMsg *LogMessage) {
	var writer io.Writer

	switch logMsg.Sev {
	case Off:
		break
	case ErrLvl:
		writer = cw.errorWriter
	case WarnLvl:
		writer = cw.warnWriter
	case InfoLvl:
		writer = cw.infoWriter
	case DebugLvl:
		writer = cw.debugWriter
	case TraceLvl:
		writer = cw.traceWriter
	}
	if writer != nil {
		writeLogMsg(writer, logMsg)
	}
}

// Close close consolewriter
func (cw *ConsoleWriter) Close() error {
	return nil
}
