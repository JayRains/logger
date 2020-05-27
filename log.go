package logger

import (
	"encoding/json"
	"fmt"
	"github.com/eliot-jay/logger/model"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
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
	  "log": {
		"file_name": "app.log",
		"file_cording": true,
		"level": "DEBUG",
		"identifier": "$",
		"time_format": "2006-01-02 15:04:05"
	  }
	}


*/

type config interface {
	config(string) (*model.Logger, error)
	decode(interface{}, reflect.Type, reflect.Value)
}

type logInfo struct {
	when      string
	path      string
	msg       string
	intactLog string
}

type logger struct {
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

func (l *logger) config(path string) (*model.Logger, error) {
	defer func() {
		fatal := recover()
		if fatal != nil {
			l.SERIOUS(err)
		}
	}()
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//start decode json config
	var decode map[string]interface{}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &decode)
	if err != nil {
		return nil, err
	}

	var set model.Logger
	l.decode(decode, reflect.TypeOf(&set), reflect.ValueOf(&set))
	return &set, nil
}

func (l *logger) decode(obj interface{}, rType reflect.Type, rValue reflect.Value) {
	rObj := reflect.ValueOf(obj)
	OneKey := rObj.MapKeys()

	for _, k := range OneKey {
		TowKey := rObj.MapIndex(k).Interface().(map[string]interface{})

		rvalue1 := reflect.ValueOf(TowKey)
		//take two key
		OkTowKey := rvalue1.MapKeys()

		for _, k1 := range OkTowKey {
			for i := 0; i < rValue.Elem().NumField(); i++ {
				if rType.Elem().Field(i).Tag.Get("json") == k1.String() {

					switch rType.Elem().Field(i).Type.Kind() {
					case reflect.String:
						rValue.Elem().Field(i).SetString(rvalue1.MapIndex(k1).Interface().(string))
					case reflect.Bool:
						rValue.Elem().Field(i).SetBool(rvalue1.MapIndex(k1).Interface().(bool))
					case reflect.Int64:
						rValue.Elem().Field(i).SetInt(int64(rvalue1.MapIndex(k1).Interface().(float64)))
					}

				}
			}

		}
	}
}

/*
check log level . current configure  level greater than  print level then print
*/
func (l *logger) state(level string, f interface{}, v ...interface{}) {
	if levelInt[l.level] >= levelInt[level] {
		l.lock.Lock()
		l.msg = formatLog(f, v...)
		l.handleText(level)
		l.lock.Unlock()
	}
}

func (l *logger) handleText(level string) {
	l.when = l.nowTime(level)
	l.path = initPrint()
	l.intactLogger()
	if l.retrieveFunc != nil {
		l.retrieveFunc(l.intactLog)
	}
	// open file cording
	if l.fileCording {
		writeDisk <- l.intactLog
	}
	if l.color {
		l.printRow(l.handlerColor(level))
	} else {
		l.printRow(strings.Split(l.intactLog,"\n")[0])
	}

}

func (l *logger) handlerColor(level string) string {
	return fmt.Sprintf("%v %v %v", colors[levelInt[level]](l.when), colors[levelInt[underline]](l.path), l.identifier+": "+colors[levelInt[level]](l.msg))
}

func (l *logger) nowTime(level string) string {
	return fmt.Sprintf("%v [%s]", time.Now().Format(l.timeFormat), level)
}

func (l *logger) printRow(msg string) {
	if _, err := os.Stdout.Write(append([]byte(msg), '\n')); err != nil {
		l.ERROR(err)
	}
}

func (l *logger) intactLogger() {
	l.intactLog = fmt.Sprintf("%v %v %v\n", l.when, l.path, l.identifier+": "+l.msg)
}

func (l *logger) ReceiveLog(handle handlerLog) {
	l.retrieveFunc = handle
}

func (l *logger) writeFile() {
	if l.fileCording {
		go func() {
			write, err := os.OpenFile(l.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				l.SERIOUS("Creat log file Failed")
				return
			}
			defer write.Close()
			for {
				w := <-writeDisk
				_, _ = write.WriteString(w)
			}
		}()
	}
}

type handlerLog func(log string)

func NewLogger(level, identifier, timeFormat, savePath string, fileCording bool, color bool) *logger {
	log := &logger{
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

func NewLogByJsonFile(JsonPath string) (*logger, error) {
	old := logger{}
	set, _ := old.config(JsonPath)
	return NewLogger(
		set.Level,
		set.Identifier,
		set.TimeFormat,
		set.SavePath,
		set.FileCording,
		set.Color,
	), nil

}

// windows system please close color
func DefaultLogger(writeFile bool, savePath string, color bool) *logger {
	return NewLogger(
		"DEBUG",
		"MSG",
		"2006-01-02 15:04:05",
		savePath,
		writeFile,
		color,
	)

}
