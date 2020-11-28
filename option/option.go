package option

//LoggerLogOption is option option
type LogOption func(*Logger)

//close the row log write to disk
func OffWrite() LogOption {
	return func(l *Logger) {
		l.OnWrite = false
		l.option = append(l.option, CallOffOnWrite)
	}
}

//close the row log color to console
func OffColor() LogOption {
	return func(l *Logger) {
		l.OnColor = false
		l.option = append(l.option, CallOffOnColor)
	}
}

//close the row log print to console
func OffConsole() LogOption {
	return func(l *Logger) {
		l.OnConsole = false
		l.option = append(l.option, CallOffOnConsole)
	}
}

// enable sprint console log
func OnConsole() LogOption {
	return func(l *Logger) {
		l.OnConsole = true

	}
}
