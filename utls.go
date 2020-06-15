package logger

import (
	"encoding/json"
	"fmt"
	"github.com/eliot-jay/logger/model"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)


func (l *Logger) config(path string) (*model.Logger, error) {
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

func (l *Logger) decode(obj interface{}, rType reflect.Type, rValue reflect.Value) {
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
func (l *Logger) state(level string, f interface{}, v ...interface{}) {
	if levelInt[l.level] >= levelInt[level] {
		l.lock.Lock()
		l.msg = formatLog(f, v...)
		l.handleText(level)
		l.lock.Unlock()
	}
}

func (l *Logger) handleText(level string) {
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

func (l *Logger) handlerColor(level string) string {
	return fmt.Sprintf("%v %v %v", colors[levelInt[level]](l.when), colors[levelInt[underline]](l.path), l.identifier+": "+colors[levelInt[level]](l.msg))
}

func (l *Logger) nowTime(level string) string {
	return fmt.Sprintf("%v [%s]", time.Now().Format(l.timeFormat), level)
}

func (l *Logger) printRow(msg string) {
	if _, err := os.Stdout.Write(append([]byte(msg), '\n')); err != nil {
		l.ERROR(err)
	}
}

func (l *Logger) intactLogger() {
	l.intactLog = fmt.Sprintf("%v %v %v\n", l.when, l.path, l.identifier+": "+l.msg)
}

func (l *Logger) WithLogMiddleware(handle handlerLog) {
	l.retrieveFunc = handle
}

func (l *Logger) writeFile() {
	if l.fileCording {
		go func() {
			write, err := os.OpenFile(l.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				l.SERIOUS("Creat log file Failed")
				return
			}
			defer write.Close()
			for {
				_, _ = write.WriteString(<-writeDisk)
			}
		}()
	}
}




const (
	//level field
	debug   = "DBUG"
	info    = "INFO"
	warn    = "WARN"
	err     = "ERRO"
	serious = "SERI"

	//system variable
	skip      = 5
	underline = "Underline"
)

var (
	colors = []brush{
		newBrush("1;41"), // Fatal              红色底	0
		newBrush("1;31"), // Error              红色		1
		newBrush("1;45"), // Warn               紫红底	2
		newBrush("1;34"), // Info               蓝色		3
		newBrush("1;32"), // Debug              绿色		4
		newBrush("4;36"), //underline          青色+下划线 5
	}
	levelInt  = make(map[string]int)
	writeDisk = make(chan string)
)

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//auto jump sprint
		} else {
			//add format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		//add format char
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}

func ToString(i int) string {
	return ":" + strconv.Itoa(i)
}

//initialize log level
func init() {
	//log level
	levelInt[debug] = 4
	levelInt[info] = 3
	levelInt[warn] = 2
	levelInt[err] = 1
	levelInt[serious] = 0

	//file path color
	levelInt[underline] = 5
}

type brush func(string) string

//color brush
func newBrush(color string) brush {
	pre := "\033[" // \033[ 1; 32m%s  \033[0m
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

//initialize log got print Row number
func initPrint() string {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(skip, rpc[:])
	if n < 1 {
		return ""
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	currentDir, _ := os.Getwd()
	return strings.Replace(frame.File, currentDir+"/", "", -1) + ToString(frame.Line)
}



