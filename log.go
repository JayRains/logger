package main


import (
	"fmt"
	"os"
	"sync"
	"time"
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
	  "logger": {
		"file_name": "app.log",
		"file_cording": true,
		"level": "DEBUG",
		"identifier": "$",
		"time_format": "2006-01-02 15:04:05"
	  }
	}


*/

type logInfo struct {
	when string
	path string
	msg  string
	intactLog  string
}

type logger struct {
	lock         *sync.Mutex
	fileName     string
	level        string
	identifier   string
	timeFormat   string
	fileCording  bool
	retrieveFunc handlerLog
	logInfo
}

func (l *logger) DEBUG(f interface{}, v ...interface{}) {
	l.state(debug, f, v...)
}

func (l *logger) INFO(f interface{}, v ...interface{}) {
	l.state(info, f, v...)
}

func (l *logger) WARN(f interface{}, v ...interface{}) {
	l.state(warn, f, v...)
}

func (l *logger) ERROR(f interface{}, v ...interface{}) {
	l.state(err, f, v...)
}

func (l *logger) SERIOUS(f interface{}, v ...interface{}) {
	l.state(serious, f, v...)
}

/*
check log level . current configure  level greater than  print level then print
*/
func (l *logger) state(level string, f interface{}, v ...interface{}) {
	if levelInt[l.level] >= levelInt[level] {
		l.msg = formatLog(f, v...)
		l.handleText(level)
	}
}

func (l *logger) handleText(level string) {
	l.lock.Lock()
	l.when = l.nowTime(level)
	l.path = initPrint()
	l.intactLogger()
	if l.retrieveFunc!=nil{
		l.retrieveFunc(l.intactLog)
	}
	// open file cording
	if l.fileCording {log <- l.intactLog}
	l.printRow(l.handlerColor(level))
	l.lock.Unlock()
}

func (l *logger) handlerColor(level string) string {
	return fmt.Sprintf("%v %v %v", colors[levelInt[level]](l.when), colors[levelInt[underline]](l.path), l.identifier+": "+colors[levelInt[level]](l.msg))
}

func (l *logger) nowTime(level string) string {
	return fmt.Sprintf("%v [%s]", time.Now().Format(l.timeFormat), level)
}

func (l *logger) printRow(msg string) {
	if _, err := os.Stdout.Write(append([]byte(msg), '\n')); err != nil {l.ERROR(err)}
}

func (l *logger) intactLogger(){
	l.intactLog = fmt.Sprintf("%v %v %v\n", l.when, l.path, l.identifier+": "+l.msg)
}

func (l *logger)ReceiveLog(handle handlerLog)  {
	l.retrieveFunc = handle
}

func (l *logger) writeFile() {
	if l.fileCording {
		go func() {
			write, err := os.OpenFile("./"+l.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {l.SERIOUS("Creat logger file Failed");return}
			defer write.Close()
			for {msg := <-log;_, _ = write.WriteString(msg)}
		}()
	}
}

type handlerLog func(log string)

func NewLogger(level, identifier, timeFormat, filename string, fileCording bool) *logger {
	 log := &logger{
		lock:        &sync.Mutex{},
		level:       level,
		identifier:  identifier,
		timeFormat:  timeFormat,
		fileCording: fileCording,
		fileName:    filename,
	}
	if log.fileCording {log.writeFile()}
	return log
}

func NewLogByJsonFile(path string) (*logger, error) {
	if set, err := newConfig(path); err != nil {return nil, err} else {
		log := &logger{
			lock:        &sync.Mutex{},
			level:       set.Level,
			identifier:  set.Identifier,
			timeFormat:  set.TimeFormat,
			fileCording: set.FileCording,
			fileName:    set.FileName,
		};log.writeFile();return log, nil}
}

func DefaultLogger(writeFile bool) *logger {
	return NewLogger("DEBUG", "MSG", "2006-01-02 15:04:05", "app.log", writeFile)

}
