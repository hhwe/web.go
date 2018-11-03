package main

import (
	"log"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	// Mock a app as http Handler
	logger := NewLogger(os.Stderr, log.LstdFlags)
	logger.SetLever("debug")
	logger.Debug("the debug message!")
	log.Println("the")
}
