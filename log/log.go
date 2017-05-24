package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
)

type LogLevel uint8

const (
	TRACE = LogLevel(iota)
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	logLevel             LogLevel
	category             string
	format               string
	isStackTraceRequired bool
	mu                   sync.Mutex
}

// NewLogger returns a new Logger with the given configuration
func NewLogger(logLevel LogLevel, category, format string) Logger {
	return Logger{logLevel, category, format, true, sync.Mutex{}}
}

// Trace logs the given message in TRACE level
func (l *Logger) Trace(message string) {
	if l.logLevel <= TRACE {
		l.log(message, TRACE)
	}
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	if l.logLevel <= TRACE {
		l.log(fmt.Sprintf(format, args...), TRACE)
	}
}

// Debug logs the given message in DEBUG level
func (l *Logger) Debug(message string) {
	if l.logLevel <= DEBUG {
		l.log(message, DEBUG)
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.logLevel <= DEBUG {
		l.log(fmt.Sprintf(format, args...), DEBUG)
	}
}

// Info logs the given message in INFO level
func (l *Logger) Info(message string) {
	if l.logLevel <= INFO {
		l.log(message, INFO)
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.logLevel <= INFO {
		l.log(fmt.Sprintf(format, args...), INFO)
	}
}

// Warn logs the given message in WARN level
func (l *Logger) Warn(message string) {
	if l.logLevel <= WARN {
		l.log(message, WARN)
	}
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.logLevel <= WARN {
		l.log(fmt.Sprintf(format+"\n", args...), WARN)
	}
}

// Error logs the given message in ERROR level
func (l *Logger) Error(message string) {
	if l.logLevel <= ERROR {
		l.log(message, ERROR)
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.logLevel <= ERROR {
		l.log(fmt.Sprintf(format, args...), ERROR)
	}
}

// Fatal logs the given message in FATAL level
func (l *Logger) Fatal(message string) {
	if l.logLevel <= FATAL {
		l.log(message, FATAL)
	}
}

func (l *Logger) FatalF(format string, args ...interface{}) {
	if l.logLevel <= FATAL {
		l.log(fmt.Sprintf(format, args...), FATAL)
	}
}

// log outputs the given message in the given log level to the standard output
func (l *Logger) log(message string, level LogLevel) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	file = filepath.Base(file)
	l.mu.Lock()
	defer l.mu.Unlock()
	var levelStr string
	switch level {
	case TRACE:
		levelStr = "TRACE"
	case DEBUG:
		levelStr = "DEBUG"
	case INFO:
		levelStr = "INFO"
	case WARN:
		levelStr = "WARN"
	case ERROR:
		levelStr = "ERROR"
	case FATAL:
		levelStr = "FATAL"
	default:
		levelStr = "LOG"
	}
	fmt.Printf("%s:%3d %5s - %s\n", file, line, levelStr, message)
}
