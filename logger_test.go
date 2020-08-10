package logger

import (
	"fmt"
	"github.com/eliot-jay/logger/logger"
	"github.com/eliot-jay/logger/public"
	"testing"
)

func TestNewLogger(t *testing.T) {
	log, _ := logger.NewLoggerBy("./conf/logger.yaml")
	defer log.Destroy()
	log.Debug("debug")
	fmt.Println(log.SPrintf(public.GenTraceID(), "warn", "hello world"))
	log.Info("Info")
	log.Warn("Warn")
	// close the row print color
	log.Error("Error")
	// close the row write to disk
	log.Serious("hello world")

}
