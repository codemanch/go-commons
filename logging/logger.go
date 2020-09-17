package logging

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/appmanch/go-commons/textutils"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/appmanch/go-commons/fsutils"
	"github.com/appmanch/go-commons/misc"
)

//Severity of the logging levels
type Severity int

//LogConfig - Configuration & Settings for the logger.
type LogConfig struct {

	//Format of the log. valid values are text,json
	//Default is text
	Format string `json:"format,omitempty" yaml:"format,omitempty"`
	//Async Flag to indicate if the writing of the flag is asynchronous.
	//Default value is false
	Async bool `json:"async,omitempty" yaml:"async,omitempty"`
	//QueueSize to indicate the number log routines that can be queued  to use in background
	//This value is used only if the async value is set to true.
	//Default value for the number items to be in queue 512
	QueueSize int `json:"queueSize,omitempty" yaml:"numRoutines,omitempty"`
	//Date - Defaults to  time.RFC3339 pattern
	DatePattern string `json:"datePattern,omitempty" yaml:"datePattern,omitempty"`
	//IncludeFunction will include the calling function name  in the log entries
	//Default value : false
	IncludeFunction bool `json:"includeFunction,omitempty" yaml:"includeFunction,omitempty"`
	//IncludeLineNum ,includes Line number for the log file
	//If IncludeFunction Line is set to false this config is ignored
	IncludeLineNum bool `json:"includeLineNum,omitempty" yaml:"includeLineNum,omitempty"`
	//DefaultLvl that will be used as default
	DefaultLvl string `json:"defaultLvl" yaml:"defaultLvl"`
	//PackageConfig that can be used to
	PkgConfigs []*PackageConfig `json:"pkgConfigs" yaml:"pkgConfigs"`
	//Writers writers for the logger. Need one for all levels
	//If a writer is not found for a specific level it will fallback to os.Stdout if the level is greater then Warn and os.Stderr otherwise
	Writers []*WriterConfig `json:"writers" yaml:"writers"`
}

// PackageConfig configuration
type PackageConfig struct {
	//PackageName
	PackageName string `json:"pkgName" yaml:"pkgName"`
	//Level to be set valid values : OFF,ERROR,WARN,INFO,DEBUG,TRACE
	Level string `json:"level" yaml:"level"`
}

//WriterConfig struct
type WriterConfig struct {
	//File reference. Non mandatory but one of file or console logger is required.
	File *FileConfig `json:"file,omitempty" yaml:"file,omitempty"`
	//Console reference
	Console *ConsoleConfig `json:"console,omitempty" yaml:"console,omitempty"`
}

//FileConfig - Configuration of file based logging
type FileConfig struct {
	//FilePath for the file based log writer
	DefaultPath string `json:"defaultPath" yaml:"defaultPath"`
	ErrorPath   string `json:"errorPath" yaml:"errorPath"`
	WarnPath    string `json:"warnPath" yaml:"warnPath"`
	InfoPath    string `json:"infoPath" yaml:"infoPath"`
	DebugPath   string `json:"debugPath" yaml:"debugPath"`
	TracePath   string `json:"tracePath" yaml:"tracePath"`
}

// ConsoleConfig - Configuration of console based logging. All Log Levels except ERROR and WARN are written to os.Stdout
// The ERROR and WARN log levels can be written  to os.Stdout or os.Stderr, By default they go to os.Stderr
type ConsoleConfig struct {
	//WriteErrToStdOut write error messages to os.Stdout .
	WriteErrToStdOut bool `json:"errToStdOut" yaml:"errToStdOut"`
	//WriteWarnToStdOut write warn messages to os.Stdout .
	WriteWarnToStdOut bool `json:"warnToStdOut" yaml:"warnToStdOut"`
}

type LogMessage struct {
	Time   time.Time `json:"timestamp"`
	FnName string    `json:"function,omitempty" `
	Line   int       `json:"line,omitempty"`
	Msg    string    `json:"msg"`
	Sev    Severity  `json:"sev"`
}

//Logger struct.
type Logger struct {
	sev             Severity
	pkgName         string
	errorEnabled    bool
	warnEnabled     bool
	infoEnabled     bool
	debugEnabled    bool
	traceEnabled    bool
	includeFunction bool
	includeLine     bool
}

type LogWriter interface {
	InitConfig(w *WriterConfig)
	DoLog(logMsg *LogMessage)
	io.Closer
}

//Map to hold loggers. This is updated in case the log config is reloaded
var loggers = make(map[string]*Logger)

// Writers can be multiple writers
var writers []LogWriter

// Log configuration
var logConfig *LogConfig

//channel of type log message
var logMsgChannel chan *LogMessage

//Levels of the logging by severity
var Levels = [...]string{
	"OFF",
	"ERROR",
	"WARN",
	"INFO",
	"DEBUG",
	"TRACE",
}

//LevelsMap of the logging by severity string severity type
var LevelsMap = map[string]Severity{
	"OFF":   Off,
	"ERROR": ErrLvl,
	"WARN":  WarnLvl,
	"INFO":  InfoLvl,
	"DEBUG": DebugLvl,
	"TRACE": TraceLvl,
}

const (
	//Off - No logging
	Off Severity = iota
	//ErrLvl - logging only for error level.
	ErrLvl
	//WarnLvl - logging turned on for warning & error levels
	WarnLvl
	//InfoLvl - logging turned on for Info, Warning and Error levels.
	InfoLvl
	//DebugLvl - Logging turned on for Debug, Info, Warning and Error levels.
	DebugLvl
	//TraceLvl - Logging turned on for Trace,Info, Warning and error Levels.
	TraceLvl
	// LogConfigEnvProperty specifies the environment variable that would specify the file location
	LogConfigEnvProperty = "LOG_CONFIG_FILE"
	//DefaultlogFilePath specifies the location where the application should search for log config if the LogConfigEnvProperty is not specified
	DefaultlogFilePath = "./log-config.json"
)

func init() {
	Configure(loadConfig())
}

func Configure(l *LogConfig) {
	logConfig = l
	if l.DatePattern == "" {
		l.DatePattern = time.RFC3339
	}
	if l.Async {

		if l.QueueSize == 0 {
			l.QueueSize = 512
		}
		logMsgChannel = make(chan *LogMessage, l.QueueSize)
		go doAsyncLog()
	}
	if l.Writers != nil {
		for _, w := range l.Writers {
			if w.File != nil {
				fw := &FileWriter{}
				fw.InitConfig(w)
				writers = append(writers, fw)
			} else if w.Console != nil {
				cw := &ConsoleWriter{}
				cw.InitConfig(w)
				writers = append(writers, cw)
			}

		}
	}
}

//Update the flags based on the severity level
func (l *Logger) updateLvlFlags() error {
	if l.sev < 0 || l.sev > 5 {
		return errors.New("Invalid severity ")
	}
	l.errorEnabled = l.sev >= 1
	l.warnEnabled = l.sev >= 2
	l.infoEnabled = l.sev >= 3
	l.debugEnabled = l.sev >= 4
	l.traceEnabled = l.sev == 5
	return nil
}

//loadDefaultConfig function with load the default configuration
func loadDefaultConfig() *LogConfig {

	return &LogConfig{
		Format:      "text",
		Async:       false,
		DatePattern: time.RFC3339,
		DefaultLvl:  Levels[InfoLvl],
		Writers: []*WriterConfig{
			{
				Console: &ConsoleConfig{},
			},
		},
	}
}

//loadConfig function will load the log configuration.
func loadConfig() *LogConfig {
	var logConfig = &LogConfig{}
	fileName := misc.GetEnvAsString(LogConfigEnvProperty, DefaultlogFilePath)
	if fsutils.FileExists(fileName) {
		contentType := fsutils.LookupContentType(fileName)
		if contentType == "application/json" {
			logConfigFile, err := os.Open(fileName)
			if err != nil {
				writeLog(os.Stderr, "Unable to open the log config file using default log configuration", err)
				logConfig = loadDefaultConfig()
			} else {
				defer logConfigFile.Close()
				bytes, _ := ioutil.ReadAll(logConfigFile)
				err = json.Unmarshal(bytes, &logConfig)
				if err != nil {
					writeLog(os.Stderr, "Unable to open the log config file using default log config", err)
					logConfig = loadDefaultConfig()
				}
			}
		} else {
			writeLog(os.Stderr, "Invalid file format supported format : application/json . Loading Default configuration")
			logConfig = loadDefaultConfig()
		}
		//TODO Add yaml support once its available
	} else {
		writeLog(os.Stderr, "Log Config file not found. Loading default configuration")
		logConfig = loadDefaultConfig()
	}
	return logConfig
}

//GetLogger function will return the logger for that package
func GetLogger() *Logger {
	pc, _, _, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	fnNameSplit := strings.Split(details.Name(), "/")
	pkgFnName := strings.Split(fnNameSplit[len(fnNameSplit)-1], ".")
	pkgName := pkgFnName[0]

	if logger, ok := loggers[pkgName]; ok {
		return logger
	}
	Level := logConfig.DefaultLvl

	if logConfig.PkgConfigs != nil && len(logConfig.PkgConfigs) > 0 {
		for _, pkgConfig := range logConfig.PkgConfigs {
			if pkgConfig.PackageName == pkgName {
				Level = pkgConfig.Level
			}
		}
	}

	logger := &Logger{
		sev:             LevelsMap[Level],
		pkgName:         pkgName,
		includeFunction: logConfig.IncludeFunction,
		includeLine:     logConfig.IncludeLineNum,
	}
	_ = logger.updateLvlFlags()
	loggers[pkgName] = logger
	return logger
}

func writeLogMsg(writer io.Writer, logMsg *LogMessage) {
	if logConfig.Format == "json" {
		data, _ := json.Marshal(logMsg)
		//TODO check if there is a better way
		_, _ = writer.Write(data)
		writeLog(writer, "")
	} else if logConfig.Format == "text" {
		if logMsg.FnName != textutils.EmptyStr {
			writeLog(writer, logMsg.Time.Format(logConfig.DatePattern), Levels[logMsg.Sev], logMsg.FnName+":"+strconv.Itoa(logMsg.Line), logMsg.Msg)
		} else {
			writeLog(writer, logMsg.Time.Format(logConfig.DatePattern), Levels[logMsg.Sev], logMsg.Msg)
		}

	}
}

//createLogMessage function creates a new log message with actual content variables
func handleLog(sev Severity, l *Logger, msg string) {

	logMsg := &LogMessage{
		Time: time.Now(),
		Msg:  msg,
		Sev:  sev,
	}

	if l.includeFunction {
		pc, _, no, _ := runtime.Caller(2)
		details := runtime.FuncForPC(pc)
		fnNameSplit := strings.Split(details.Name(), "/")
		logMsg.FnName = fnNameSplit[len(fnNameSplit)-1]

		if l.includeLine {
			logMsg.Line = no
		}
	}
	if logConfig.Async {
		logMsgChannel <- logMsg

	} else {
		for _, w := range writers {
			w.DoLog(logMsg)
		}
	}
}

func doAsyncLog() {

	for logMsg := range logMsgChannel {
		for _, w := range writers {
			w.DoLog(logMsg)
		}
	}

}

//writeLog will write to the io.Writer interface
func writeLog(w io.Writer, a ...interface{}) {
	fmt.Fprintln(w, a...)
}

//String method to get the Severity String
func (sev Severity) String() (string, error) {
	if sev < 0 || sev > 5 {
		return "", errors.New("Invalid severity ")
	}
	return Levels[sev], nil
}

//IsEnabled function returns if the current
func (l *Logger) IsEnabled(sev Severity) bool {
	return sev <= TraceLvl && sev >= l.sev
}

//Error Logger
func (l *Logger) Error(a ...interface{}) {
	if l.errorEnabled && a != nil && len(a) > 0 {
		handleLog(ErrLvl, l, fmt.Sprint(a...))
	}
}

//ErrorF Logger with formatting of the messages
func (l *Logger) ErrorF(f string, a ...interface{}) {
	if l.errorEnabled {
		handleLog(ErrLvl, l, fmt.Sprintf(f, a...))
	}
}

//Warn Logger
func (l *Logger) Warn(a ...interface{}) {
	if l.warnEnabled && a != nil && len(a) > 0 {
		handleLog(WarnLvl, l, fmt.Sprint(a...))
	}
}

//WarnF Logger with formatting of the messages
func (l *Logger) WarnF(f string, a ...interface{}) {
	if l.warnEnabled {
		handleLog(WarnLvl, l, fmt.Sprintf(f, a...))

	}
}

//Info Logger
func (l *Logger) Info(a ...interface{}) {
	if l.infoEnabled && a != nil && len(a) > 0 {
		handleLog(InfoLvl, l, fmt.Sprint(a...))

	}
}

//InfoF Logger
func (l *Logger) InfoF(f string, a ...interface{}) {
	if l.infoEnabled {
		handleLog(InfoLvl, l, fmt.Sprintf(f, a...))

	}
}

//Debug Logger
func (l *Logger) Debug(a ...interface{}) {
	if l.debugEnabled && a != nil && len(a) > 0 {
		handleLog(DebugLvl, l, fmt.Sprint(a...))
	}
}

//DebugF Logger
func (l *Logger) DebugF(f string, a ...interface{}) {
	if l.debugEnabled {
		handleLog(DebugLvl, l, fmt.Sprintf(f, a...))
	}
}

//Trace Logger
func (l *Logger) Trace(a ...interface{}) {
	if l.traceEnabled && a != nil && len(a) > 0 {
		handleLog(TraceLvl, l, fmt.Sprint(a...))

	}
}

//TraceF Logger
func (l *Logger) TraceF(f string, a ...interface{}) {
	if l.traceEnabled {
		handleLog(TraceLvl, l, fmt.Sprintf(f, a...))
	}
}
