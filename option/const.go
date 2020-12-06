package option

const (
	DEBUG   = "DEBUG"
	INFO    = "IN-FO"
	WARN    = "WA-RN"
	ERRNO   = "ERROR"
	SERIOUS = "SERIO"
	FATAL   = "FATAL"
	Linux   = "linux"
	Skip    = 4

	Blue   = "#1976D2"
	Green  = "#2ecc71"
	Red    = "#e74c3c"
	Purple = "#9b59b6"
)

const Trace = "|| Track"

const (
	CallOffOnWrite = iota
	CallOffOnColor
	CallOffOnConsole
)
