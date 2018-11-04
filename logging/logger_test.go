package logging

import (
	"log"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := NewLogger(os.Stderr, log.LstdFlags)
	logger.SetLever("info")

	logger.Debug("debug message!")
	logger.Info("info message!")
	logger.Warning("warnig message!")
	logger.Error("error message!")
}
