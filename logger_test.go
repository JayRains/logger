package logger

import (
	"github.com/eliot-jay/logger/public"
	"github.com/eliot-jay/logger/register"
	"testing"
)
var (
	filter =  make(map[string]bool)
	count = 0
)

func TestNewLogger(t *testing.T) {
	log ,err := register.NewLogger("./config.yaml")
	if err!=nil{
		panic(err)
	}
	defer log.Destroy()
	log.Debug("hello world")
	log.Info("Info")
	log.Warn("Warn")
	log.Error("Error")
	log.Serious("Serious")
}


func BenchmarkAll(b *testing.B) {
	for i:=0;i<b.N;i++{
		iterateFilter(public.GenTraceID())
	}
}
func iterateFilter(hash string ) bool {
	if filter[hash]{
		count++
		return true
	}
	filter[hash] = true
	return false
}