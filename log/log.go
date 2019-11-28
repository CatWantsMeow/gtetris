package log

import (
	"fmt"
	"strings"
	"time"
)

const (
	LevelDebug   = "D"
	LevelInfo    = "I"
	LevelWarning = "W"
	LevelError   = "E"
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
	logger = Logger{}
)

func Debug(msg string, args ...interface{}) {
	logger.Log(LevelDebug, msg, args...)
}

func Info(msg string, args ...interface{}) {
	logger.Log(LevelInfo, msg, args...)
}

func Warning(msg string, args ...interface{}) {
	logger.Log(LevelWarning, msg, args...)
}

func Error(msg string, args ...interface{}) {
	logger.Log(LevelError, msg, args...)
}

func String(tail, width int) string {
	return logger.String(tail, width)
}
