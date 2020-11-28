package logger

import (
	"github.com/eliot-jay/logger/option"
	"testing"
)

func TestNewLogger(t *testing.T) {
	log := option.Default()
	println := log.Register("HELLO", option.Purple, 4)
	println("hello world")
	//log.Debug("debug")
	//log.Info("Info")
	//log.Warn("Warn")
	//// close the row print color
	//log.Debug(log.Error("Error",))
	//// close the row write to disk
	//log.Serious("hello world",option.OffWrite())

}
