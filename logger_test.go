package logger

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	var logger Console
	logger ,_ = NewLogByJsonFile("./log.json")
	logger.DEBUG("hello debug")
	logger.INFO("hello info")
	logger.ERROR("hello error")
	logger.WARN("hello warn")
	logger.SERIOUS("hello serious")


}
