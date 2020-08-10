package logger

import "github.com/eliot-jay/logger/public"

//LoggerHandler ...
//description:  LoggerHandler is logger option
type LoggerHandler func(*logger)

//OffWrite  ...
//description:  close the row log write to disk
func OffWrite() LoggerHandler {
	return func(l *logger) {
		l.OnWrite = false
		l.option = append(l.option, public.CallOffOnWrite)
	}
}

//OffColor  ...
//description:  close the row log color to console
func OffColor() LoggerHandler {
	return func(l *logger) {
		l.OnColor = false
		l.option = append(l.option, public.CallOffOnColor)
	}
}

//OffConsole  ...
//description:  close the row log print to console
func OffConsole() LoggerHandler {
	return func(l *logger) {
		l.OnConsole = false
		l.option = append(l.option, public.CallOffOnConsole)
	}
}
