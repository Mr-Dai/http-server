package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
)

type Level uint8

const (
	TRACE = Level(iota)
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	level                Level
	category             string
	format               string
	isStackTraceRequired bool
	mu                   sync.Mutex
}

// NewLogger returns a new Logger with the given configuration
func NewLogger(level Level, category, format string) *Logger {
	return &Logger{
		level:                level,
		category:             category,
		format:               format,
		isStackTraceRequired: true,
		mu:                   sync.Mutex{}}
}

// Trace logs the given message in TRACE level to the logger
func (l *Logger) Trace(message string) {
	if l.level <= TRACE {
		l.log(message, TRACE)
	}
}

// Tracef formats according to the format specifier and logs the result message
// in TRACE level to the logger
func (l *Logger) Tracef(format string, args ...interface{}) {
	if l.level <= TRACE {
		l.log(fmt.Sprintf(format, args...), TRACE)
	}
}

// Debug logs the given message in DEBUG level to the logger
func (l *Logger) Debug(message string) {
	if l.level <= DEBUG {
		l.log(message, DEBUG)
	}
}

// Debugf formats according to the format specifier and logs the result message
// in DEBUG level to the logger
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.level <= DEBUG {
		l.log(fmt.Sprintf(format, args...), DEBUG)
	}
}

// Info logs the given message in INFO level to the logger
func (l *Logger) Info(message string) {
	if l.level <= INFO {
		l.log(message, INFO)
	}
}

// Infof formats according to the format specifier and logs the result message
// in INFO level to the logger
func (l *Logger) Infof(format string, args ...interface{}) {
	if l.level <= INFO {
		l.log(fmt.Sprintf(format, args...), INFO)
	}
}

// Warn logs the given message in WARN level to the logger
func (l *Logger) Warn(message string) {
	if l.level <= WARN {
		l.log(message, WARN)
	}
}

// Warnf formats according to the format specifier and logs the result message
// in WARN level to the logger
func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.level <= WARN {
		l.log(fmt.Sprintf(format, args...), WARN)
	}
}

// Error logs the given message in ERROR level to the logger
func (l *Logger) Error(message string) {
	if l.level <= ERROR {
		l.log(message, ERROR)
	}
}

// Errorf formats according to the format specifier and logs the result message
// in ERROR level to the logger
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.level <= ERROR {
		l.log(fmt.Sprintf(format, args...), ERROR)
	}
}

// Fatal logs the given message in FATAL level to the logger
func (l *Logger) Fatal(message string) {
	if l.level <= FATAL {
		l.log(message, FATAL)
	}
}

// Fatalf formats according to the format specifier and logs the result message
// in FATAL level to the logger
func (l *Logger) Fatalf(format string, args ...interface{}) {
	if l.level <= FATAL {
		l.log(fmt.Sprintf(format, args...), FATAL)
	}
}

// log outputs the given message in the given log level to the standard output
func (l *Logger) log(message string, level Level) {
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
