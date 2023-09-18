package logger

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, level Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: level,
	}
}

func (l *Logger) print(level Level, msg string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}
	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties"`
		Trace      string            `json:"-"`
	}{
		Level:      level.String(),
		Time:       time.Now().Format("2006-01-02"),
		Message:    msg,
		Properties: properties,
	}
	if level >= l.minLevel {
		aux.Trace = string(debug.Stack()[:30])
	}
	var err error
	var line []byte

	line, err = json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ":\tunable to marshal log\n" + err.Error())
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))

}
func (l *Logger) Write(msg []byte) (n int, err error) {
	return l.print(LevelError, string(msg), nil)
}

func (l *Logger) PrintInfo(msg string, properties map[string]string) {
	l.print(LevelInfo, msg, properties)
}

func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(-1)
}
