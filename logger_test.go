package logger

import (
	"fmt"
	"github.com/eliot-jay/logger/register"
	"testing"
)

var (
	filter = make(map[string]bool)
	count  = 0
)

func TestNewLogger(t *testing.T) {

	log := register.NewDefaultLogger()
	defer log.Destroy()
	fmt.Println(log.Sprint("warn", "hello world").Text())
	log.Debug("hello world")
	log.Info("Info")
	log.Warn("Warn")
	log.Error("Error")
	log.Serious("Serious")

}
