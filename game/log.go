package game

import (
	"fmt"
	"strings"
	"time"
)

const (
	LogLevelDebug   = "D"
	LogLevelInfo    = "I"
	LogLevelWarning = "W"
	LogLevelError   = "E"
)

type Logger struct {
	entries []string
}

func (l *Logger) Log(level string, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	timestamp := time.Now().Format("15:04:05.0000")
	entry := level + " [" + timestamp + "] " + msg
	l.entries = append(l.entries, entry)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.Log(LogLevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.Log(LogLevelInfo, msg, args...)
}

func (l *Logger) Warning(msg string, args ...interface{}) {
	l.Log(LogLevelWarning, msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.Log(LogLevelError, msg, args...)
}

func (l *Logger) String(tail int, width int) string {
	start := len(l.entries) - tail
	if start < 0 {
		start = 0
	}

	var sb strings.Builder
	for _, entry := range l.entries[start:] {
		w := len(entry)
		if w > width {
			w = width
		}

		sb.WriteString(entry[:w])
		sb.WriteString("\n")
	}
	return sb.String()
}

var (
	Log = Logger{}
)
