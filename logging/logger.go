// logging is based on golang standar package log.
// colorful, structured, leveled logging in Go.
package logging

import (
	"errors"
	"fmt"
	"log"
	"os"
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

// Logger is base on log package with log level and colorful output
type Logger struct {
	Logger *log.Logger
	level  levelType
	color  bool
}

func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stderr, "", 0),
		level:  INFO,
		color:  false,
	}
}

// SetLevel change the default level of Logger
func (l *Logger) SetLever(name string) {
	level, ok := nameLevel[strings.ToUpper(name)]
	if !ok {
		panic(LevelNotSupported)
	}
	l.level = level
}

// SetColor logged output with colorful string
func (l *Logger) SetColor(color bool) {
	l.color = color
}

// color terminal console output colorful string
func color(color uint8, s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", color, s)
}

func (l *Logger) output(level levelType, msg string) {
	if level >= l.level {
		if l.color {
			l.Logger.Println(color(94-uint8(level), levelName[level]) + " " + msg)
		} else {
			l.Logger.Println(levelName[level] + " " + msg)
		}
	}
}

// we can use implete Debugf(format string, msd ...interface{})
// Logger.Info(fmt.Sprintf("%v", time.Since(start)))
func (l *Logger) Debug(msg ...interface{}) {
	l.output(DEBUG, fmt.Sprint(msg...))
}

func (l *Logger) Info(msg ...interface{}) {
	l.output(INFO, fmt.Sprint(msg...))
}

func (l *Logger) Warning(msg ...interface{}) {
	l.output(WARNING, fmt.Sprint(msg...))
}

func (l *Logger) Error(msg ...interface{}) {
	l.output(ERROR, fmt.Sprint(msg...))
}
