package util

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// LogLevel represents the verbosity level for logging.
type LogLevel int

const (
	// LevelQuiet suppresses all output except errors.
	LevelQuiet LogLevel = 0
	// LevelNormal shows standard output.
	LevelNormal LogLevel = 1
	// LevelVerbose shows additional information.
	LevelVerbose LogLevel = 2
	// LevelDebug shows debug information.
	LevelDebug LogLevel = 3
)

// Logger provides configurable logging with verbosity levels.
type Logger struct {
	level  LogLevel
	out    io.Writer
	errOut io.Writer
	mu     sync.Mutex
}

var defaultLogger = NewLogger(LevelNormal)

// NewLogger creates a new logger with the specified verbosity level.
func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level:  level,
		out:    os.Stdout,
		errOut: os.Stderr,
	}
}

// SetLevel sets the verbosity level.
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetOutput sets the output writer.
func (l *Logger) SetOutput(out io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = out
}

// SetErrorOutput sets the error output writer.
func (l *Logger) SetErrorOutput(out io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.errOut = out
}

// Error prints error messages (always shown).
func (l *Logger) Error(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.errOut, "❌ "+format+"\n", args...)
}

// Warn prints warning messages (shown at LevelNormal and above).
func (l *Logger) Warn(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelNormal {
		fmt.Fprintf(l.out, "⚠️ "+format+"\n", args...)
	}
}

// Info prints info messages (shown at LevelNormal and above).
func (l *Logger) Info(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelNormal {
		fmt.Fprintf(l.out, format+"\n", args...)
	}
}

// Success prints success messages (shown at LevelNormal and above).
func (l *Logger) Success(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelNormal {
		fmt.Fprintf(l.out, "✅ "+format+"\n", args...)
	}
}

// Verbose prints verbose messages (shown at LevelVerbose and above).
func (l *Logger) Verbose(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelVerbose {
		fmt.Fprintf(l.out, "   "+format+"\n", args...)
	}
}

// Debug prints debug messages (shown at LevelDebug only).
func (l *Logger) Debug(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelDebug {
		fmt.Fprintf(l.out, "🔍 "+format+"\n", args...)
	}
}

// IsDebug returns true when debug logging is enabled.
func (l *Logger) IsDebug() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level >= LevelDebug
}

// Default returns the default logger instance.
func Default() *Logger {
	return defaultLogger
}
