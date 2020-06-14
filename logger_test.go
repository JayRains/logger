package logger

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	var logger *Logger
	logger = DefaultLogger(true)
	logger.DEBUG("hello debug")
	logger.INFO("hello info")
	logger.ERROR("hello error")
	logger.WARN("hello warn")
	logger.SERIOUS("hello serious")


}
