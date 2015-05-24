package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"
)

// Logger stuct contains the log details
type Logger struct {
	Date   time.Time
	Level  string
	Src    Source
	Output *string
}

// Source struct contains the source details
type Source struct {
	Func string
	Line int
}

var log = new(Logger)

func (l *Logger) caller() {
	var ok bool
	var pc uintptr
	pc, _, l.Src.Line, ok = runtime.Caller(3)
	if !ok {
		l.Src.Func = "???"
		l.Src.Line = 0
	} else {
		l.Src.Func = runtime.FuncForPC(pc).Name()
	}

}

func (l *Logger) message(format string, args []interface{}) *string {
	m := fmt.Sprintf(format, args...)
	return &m
}

func (l *Logger) now() {
	l.Date = time.Now()
}

func (l *Logger) print(level string, format string, args ...interface{}) {
	l.now()
	l.caller()
	l.Output = l.message(format, args)
	l.Level = level

	data, err := json.Marshal(l)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
}

// Info function logs a message using INFO as log level
func Info(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.print("info", "%s", s)
}

// Debug function logs a message using DEBUG as log level if DEBUG environment variable is enabled
func Debug(args ...interface{}) {
	debug := os.Getenv("DEBUG")
	if debug != "" {
		s := fmt.Sprint(args...)
		log.print("debug", "%s", s)
	}
}

// Warning function logs a message using WARNING as log level
func Warning(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.print("warning", "%s", s)
}

// Fatal function logs a message using FATAL as log level followed by a call to os.Exit(1)
func Fatal(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.print("fatal", "%s", s)
	os.Exit(1)
}

// Infof function logs a message with arguments using INFO as log level
func Infof(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.print("info", s)
}

// Debugf function logs a message with arguments using DEBUG as log level if DEBUG environment variable is enabled
func Debugf(format string, args ...interface{}) {
	debug := os.Getenv("DEBUG")
	if debug != "" {
		s := fmt.Sprintf(format, args...)
		log.print("debug", s)
	}
}

// Warningf function logs a message with arguments using WARNING as log level
func Warningf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.print("warning", s)
}

// Fatalf function logs a message with arguments using FATAL as log level followed by a call to os.Exit(1)
func Fatalf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.print("fatal", s)
	os.Exit(1)
}
