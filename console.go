package logger

import (
	"github.com/eliot-jay/logger/model"
	"reflect"
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

	json config
		{
	  "log": {
		"file_name": "app.log",
		"file_cording": true,
		"level": "DEBUG",
		"identifier": "$",
		"time_format": "2006-01-02 15:04:05"
	  }
	}


*/

type Console interface {
	config(string) (*model.Logger, error)
	decode(interface{}, reflect.Type, reflect.Value)
	DEBUG(f interface{}, v ...interface{})
	INFO(f interface{}, v ...interface{})
	WARN(f interface{}, v ...interface{})
	ERROR(f interface{}, v ...interface{})
	SERIOUS(f interface{}, v ...interface{})
}


type Logger struct {
	color        bool
	lock         *sync.Mutex
	fileName     string
	level        string
	identifier   string
	timeFormat   string
	fileCording  bool
	retrieveFunc handlerLog
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


func NewLogger(level, identifier, timeFormat, savePath string, fileCording bool, color bool) *Logger {
	log := &Logger{
		lock:        &sync.Mutex{},
		level:       level,
		identifier:  identifier,
		timeFormat:  timeFormat,
		fileCording: fileCording,
		fileName:    savePath,
		color:       color,
	}
	if log.fileCording {
		log.writeFile()
	}
	return log
}

func NewLogByJsonFile(JsonPath string) (Console, error) {
	old := Logger{}
	set, _ := old.config(JsonPath)
	return NewLogger(
		set.Level,
		set.Identifier,
		set.TimeFormat,
		set.SavePath,
		set.FileCording,
		set.OpenColor,
	), nil

}

// windows system please close color
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
	l.color = false
	l.INFO("Windows Please Off Color for False")
	return l

}
