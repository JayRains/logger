package logger

import (
	"runtime"
	"sync"
)

/*	all time format
	ANSIC       	"Mon Jan _2 15:04:05 2006"
	UnixDate		"Mon Jan _2 15:04:05 MST 2006"
	RubyDate		"Mon Jan 02 15:04:05 -0700 2006"
	RFC822			"02 Jan 06 15:04 MST"
	RFC822Z			"02 Jan 06 15:04 -0700"
	RFC850			"Monday, 02-Jan-06 15:04:05 MST"
	RFC1123			"Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z		"Mon, 02 Jan 2006 15:04:05 -0700"
	RFC3339			"2006-01-02T15:04:05Z07:00"
	RFC3339Nano		"2006-01-02T15:04:05.999999999Z07:00"
	Kitchen			"3:04PM"
	Stamp			"Jan _2 15:04:05"
	StampMilli		"Jan _2 15:04:05.000"
	StampMicro		"Jan _2 15:04:05.000000"
	StampNano		"Jan _2 15:04:05.000000000"
	RFC3339Nano1	"2006-01-02 15:04:05.999999999 -0700 MST"
	DEFAULT			"2006-01-02 15:04:05"
*/

type Console interface {
	decode(path string) Console
	DEBUG(f interface{}, v ...interface{})
	INFO(f interface{}, v ...interface{})
	WARN(f interface{}, v ...interface{})
	ERROR(f interface{}, v ...interface{})
	SERIOUS(f interface{}, v ...interface{})
}


type Logger struct {
	Color        bool `yaml:"color"`
	SavePath     string `yaml:"savepath"`
	Level        string `yaml:"level"`
	Identifier   string `yaml:"identifier"`
	TimeFormat   string `yaml:"timeformat"`
	FileCording  bool `yaml:"filecording"`


	retrieveFunc handlerLog
	lock         *sync.Mutex
	logInfo
}


type handlerLog func(log string)
type logInfo struct {
	when      string
	path      string
	msg       string
	intactLog string
}



func (l *Logger) DEBUG(f interface{}, v ...interface{}) {
	l.state(debug, f, v...)
}

func (l *Logger) INFO(f interface{}, v ...interface{}) {
	l.state(info, f, v...)
}

func (l *Logger) WARN(f interface{}, v ...interface{}) {
	l.state(warn, f, v...)
}

func (l *Logger) ERROR(f interface{}, v ...interface{}) {
	l.state(err, f, v...)
}

func (l *Logger) SERIOUS(f interface{}, v ...interface{}) {
	l.state(serious, f, v...)
}


func NewLogger(Level, Identifier, TimeFormat, savePath string, FileCording bool, Color bool) *Logger {
	log := &Logger{
		lock:        &sync.Mutex{},
		Level:       Level,
		Identifier:  Identifier,
		TimeFormat:  TimeFormat,
		FileCording: FileCording,
		SavePath:    savePath,
		Color:       Color,
	}
	log.writeFile()
	return log
}

func NewLogByConfigFile(Path string) Console {
	return NewLogger("","","","",false,false).decode(Path)
}

// windows system please close Color
func DefaultLogger()Console{
	const (
		linux = "linux"
		systemType = runtime.GOOS
	)

	l := NewLogger(
		"DBUG",
		"$",
		"2006-01-02 15:04:05",
		"",
		false,
		true,
	)

	if systemType == linux{
		return l
	}
	l.Color = false
	l.INFO("Windows Please Off Color Switch")
	return l

}
