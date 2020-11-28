package logger

import "github.com/eliot-jay/logger/option"

func NewDefault() *option.Logger {
	return option.Default()
}
func NewLoggerBy(path string) (*option.Logger, error) {
	return option.NewLoggerBy(path)
}
