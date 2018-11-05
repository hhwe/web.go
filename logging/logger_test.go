package logging

import (
	"bytes"
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestLogger(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	logger := NewLogger()
	logger.SetLever("info")
	logger.Logger.SetOutput(buf)

	// debug level can't ouptut
	logger.Debug("debug")
	expect(t, buf.String(), "")

	buf.Reset() // clear buffed stream
	logger.Info("info")
	expect(t, buf.String(), "[INFO] info\n")

	buf.Reset()
	logger.Warning("warning")
	expect(t, buf.String(), "[WARNING] warning\n")

	buf.Reset()
	logger.Error("error")
	expect(t, buf.String(), "[ERROR] error\n")
}
