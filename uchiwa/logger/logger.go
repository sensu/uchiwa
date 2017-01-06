package logger

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Logging Levels
const (
	FATAL int = iota
	WARN
	INFO
	DEBUG
	TRACE
)

var configuredLevel int

var levels = []string{"FATAL", "WARN", "INFO", "DEBUG", "TRACE"}

// Logger stuct contains the log details
type Logger struct {
	Timestamp time.Time `json:"timestamp"`
	Level     *string   `json:"level"`
	Src       *Source   `json:"src,omitempty"`
	Message   *string   `json:"message"`
}

// Source struct contains the source details
type Source struct {
	Func string `json:"func,omitempty"`
	Line int    `json:"line,omitempty"`
}

var (
		log = new(Logger)
		logMutex = &sync.Mutex{}
)

func init() {
	configuredLevel = INFO
}

func (l *Logger) caller() {
	var ok bool
	var pc uintptr
	var src Source

	pc, _, src.Line, ok = runtime.Caller(3)
	if !ok {
		src.Func = "???"
		src.Line = 0
	} else {
		src.Func = runtime.FuncForPC(pc).Name()
	}

	l.Src = &src
}

func (l *Logger) message(format string, args []interface{}) *string {
	m := fmt.Sprintf(format, args...)
	return &m
}

func (l *Logger) now() {
	l.Timestamp = time.Now()
}

func (l *Logger) print(level string, format string, args ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	l.now()
	l.Message = l.message(format, args)
	l.Level = &level

	// Are we priting logs for this level?
	if isDisabledFor(level) {
		return
	}

	// Do we need to add additional information to the message?
	if configuredLevel >= DEBUG {
		l.caller()
	}

	data, err := json.Marshal(l)
	if err != nil {
		fmt.Println(html.EscapeString(err.Error()))
		return
	}
	fmt.Println(string(data))
}

// Debug function logs a message using DEBUG as log level if DEBUG environment variable is enabled
func Debug(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.print("debug", "%s", s)
}

// Debugf function logs a message with arguments using DEBUG as log level if DEBUG environment variable is enabled
func Debugf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.print("debug", s)
}

// Fatal function logs a message using FATAL as log level followed by a call to os.Exit(1)
func Fatal(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.print("fatal", "%s", s)
	os.Exit(1)
}

// Fatalf function logs a message with arguments using FATAL as log level followed by a call to os.Exit(1)
func Fatalf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.print("fatal", s)
	os.Exit(1)
}

// Info function logs a message using INFO as log level
func Info(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.print("info", "%s", s)
}

// Infof function logs a message with arguments using INFO as log level
func Infof(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.print("info", s)
}

// Trace function logs a message using TRACE as log level
func Trace(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.print("trace", "%s", s)
}

// Tracef function logs a message with arguments using TRACE as log level
func Tracef(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.print("trace", s)
}

// Warning function logs a message using WARNING as log level
func Warning(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.print("warn", "%s", s)
}

// Warningf function logs a message with arguments using WARNING as log level
func Warningf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.print("warn", s)
}

// getLevelInt returns the integer representation of a logging level
func getLevelInt(level string) int {
	for i, name := range levels {
		if strings.EqualFold(name, level) {
			return i
		}
	}
	// Return info by default
	return INFO
}

// isDisabledFor returns true if logging is disabled for the provided level
func isDisabledFor(level string) bool {
	levelInt := getLevelInt(level)

	return levelInt > configuredLevel
}

// SetLogLevel ...
func SetLogLevel(level string) {
	configuredLevel = getLevelInt(level)
}
