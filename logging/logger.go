// logging is imitate by py standar libriry logging.py
package logging

import (
	"errors"
	"io"
	"log"
)

type levelType int

const (
	DEBUG levelType = iota
	INFO
	WARNING
	ERROR
	CRITICAL
)

var (
	LevelNotSupported = errors.New("level not supported")
)

var levelName = map[levelType]string{
	DEBUG:    "DEBUG",
	INFO:     "INFO",
	WARNING:  "WARNING",
	ERROR:    "ERROR",
	CRITICAL: "CRITICAL",
}

func GetLevelName(level levelType) string {
	return levelName[level]
}

var nameLevel = map[string]levelType{
	"DEBUG":    DEBUG,
	"INFO":     INFO,
	"WARNING":  WARNING,
	"ERROR":    ERROR,
	"CRITICAL": CRITICAL,
}

func GetNameLevel(name string) levelType {
	return nameLevel[name]
}

type logger interface {
	Output(int, string) error
	Print(...interface{})
	Printf(string, ...interface{})
}

type Logger struct {
	Level levelType
	logger
}

func NewLogger(out io.Writer, flag int) *Logger {
	return &Logger{
		logger: log.New(out, "[web.go]", flag),
		Level:  INFO,
	}
}

// SetLevel change the default level of Logger
func (l *Logger) SetLever(name string) {
	level, ok := nameLevel[name]
	if !ok {
		panic(LevelNotSupported)
	}
	l.Level = level
}

func (l *Logger) Debug(msg string) {
}

// func (l *Logger) Info(msg string) {
// }

// func (l *Logger) Warning(msg string) {
// }
