// logging is imitate by py standar libriry logging.py
package logging

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
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
	DEBUG:    "[DEBUG]",
	INFO:     "[INFO]",
	WARNING:  "[WARNING]",
	ERROR:    "[ERROR]",
	CRITICAL: "[CRITICAL]",
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
	Println(...interface{})
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

// color terminal console output colorful string
func color(color uint8, s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", color, s)
}

// SetLevel change the default level of Logger
func (l *Logger) SetLever(name string) {
	level, ok := nameLevel[strings.ToUpper(name)]
	if !ok {
		panic(LevelNotSupported)
	}
	l.Level = level
}

func (l *Logger) output(level levelType, msg string) {
	if level >= l.Level {
		l.Println(color(94-uint8(level), levelName[level]) + " " + msg)
	}
}

func (l *Logger) Debug(msg string) {
	l.output(DEBUG, msg)
}

func (l *Logger) Info(msg string) {
	l.output(INFO, msg)
}

func (l *Logger) Warning(msg string) {
	l.output(WARNING, msg)
}

func (l *Logger) Error(msg string) {
	l.output(ERROR, msg)
}
