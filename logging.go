package main

import (
	"log"
	"os"
)

// logging setting
var logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)

type Logger interface {
	debug(v ...interface{})
	info(v ...interface{})
	warning(v ...interface{})
	error(v ...interface{})
}

func NewLogger()  {
	log.New(os.Stderr, "", log.Ldate|log.Ltime)

}
