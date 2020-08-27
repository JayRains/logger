package logger

import "github.com/eliot-jay/logger/public"


//LoggerHandler is logger option
type LoggerHandler func(*logger)


//close the row log write to disk
func OffWrite() LoggerHandler {
	return func(l *logger) {
		l.OnWrite = false
		l.option = append(l.option, public.CallOffOnWrite)
	}
}


//close the row log color to console
func OffColor() LoggerHandler {
	return func(l *logger) {
		l.OnColor = false
		l.option = append(l.option, public.CallOffOnColor)
	}
}


//close the row log print to console
func OffConsole() LoggerHandler {
	return func(l *logger) {
		l.OnConsole = false
		l.option = append(l.option, public.CallOffOnConsole)
	}
}

// enable sprint console log
func OnConsole() LoggerHandler {
	return func(l *logger) {
		l.OnConsole = true
	}
}
