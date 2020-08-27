package logger

import (
	"errors"
	"fmt"
	"github.com/eliot-jay/logger/public"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Logger interface {
	Debug(f interface{}, a ...interface{}) string
	Info(f interface{}, a ...interface{}) string
	Warn(f interface{}, a ...interface{}) string
	Error(f interface{}, a ...interface{}) string
	Serious(f interface{}, a ...interface{}) string
	Fatal(f interface{}, a ...interface{})
	SPrintf(TrackID, Type string, f interface{}, a ...interface{}) string
	Destroy()
}

type LoggerConfigure struct {
	OnConsole    bool    `yaml:"on-console"`
	OnColor      bool    `yaml:"on-color"`
	OnWrite      bool    `yaml:"on-write"`
	LoggerNormal string  `yaml:"normal-file"`
	LoggerError  string  `yaml:"error-file"`
	Level        string  `yaml:"level"`
	Identifier   string  `yaml:"identifier"`
	TimeFormat   string  `yaml:"time-format"`
	MaxSize      float64 `yaml:"max-size"`
}

type logger struct {
	FilePath        string
	TraceID         string
	NowTime         string
	Msg             string
	LoggerConfigure `yaml:"loggerConfigure"`

	//private
	lock   *sync.RWMutex
	option []uint
}

func newDefaultLogger() (log *logger) {
	log = &logger{lock: &sync.RWMutex{}, option: make([]uint, 0)}
	config := `
loggerConfigure:
  on-color: true  
  on-write: false 
  on-console: true
  level: "DBUG"
  max-size: 1
  normal-file: "./logs/project.inf.1.log"
  error-file: "./logs/project.err.1.log"
  identifier: "" 
  time-format: "2006-01-02 15:04:05"
`
	log, _ = log.ShouldBind([]byte(config))
	if runtime.GOOS != public.Linux {
		log.OnColor = false
		log.Warn("Windows System Closed Color Switch")
	}
	return log
}

func (log *logger) Debug(f interface{}, a ...interface{}) string {
	defer log.lock.Unlock()
	log.inspector(public.DEBUG, f, a...)
	return log.TraceID
}

func (log *logger) Info(f interface{}, a ...interface{}) string {
	defer log.lock.Unlock()
	log.inspector(public.INFO, f, a...)
	return log.TraceID
}

func (log *logger) Warn(f interface{}, a ...interface{}) string {
	defer log.lock.Unlock()
	log.inspector(public.WARN, f, a...)
	return log.TraceID
}

func (log *logger) Error(f interface{}, a ...interface{}) string {
	defer log.lock.Unlock()
	log.inspector(public.ERRNO, f, a...)
	return log.TraceID
}

func (log *logger) Serious(f interface{}, a ...interface{}) string {
	defer log.lock.Unlock()
	log.inspector(public.SERIOUS, f, a...)
	return log.TraceID
}
func (log *logger) Fatal(f interface{}, a ...interface{}) {
	defer os.Exit(0)
	defer log.lock.Unlock()
	defer log.Destroy()
	log.inspector(public.FATAL, f, a...)
}

func (log *logger) SPrintf(TractID, Type string, f interface{}, a ...interface{}) string {
	defer log.lock.Unlock()
	if log.OnConsole {
		log.OnConsole = false
		defer func() {
			log.OnConsole = true
		}()
	}
	genTraceID = false
	log.TraceID = TractID
	log.inspector(Type, f, a...)
	genTraceID = true
	return log.Text()
}

func (log *logger) writInspector(level string) {
	if log.FileSizeInspector() {
		log.write(level)
	}
}

func (log *logger) write(level string) {
	if level == public.ERRNO || level == public.SERIOUS || level == public.FATAL {
		_, _ = public.ErrnoWrite.WriteString(fmt.Sprintf("%v %v: %v \n", log.Text(), public.Trace, log.TraceID))
	} else {
		_, _ = public.NormalWrite.WriteString(fmt.Sprintf("%v %v: %v \n", log.Text(), public.Trace, log.TraceID))
	}
}

// logger file size compare
func (log *logger) FileSizeInspector() bool {
	if !log.OnWrite {
		return false
	}
	public.Once.Do(func() {
		var err error
		public.NormalWrite, err = os.OpenFile(log.LoggerNormal, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		public.ErrnoWrite, err = os.OpenFile(log.LoggerError, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal("Open logger File failure: %v", err)
		}
	})

	var (
		err        error
		normalInfo os.FileInfo
		errnoInfo  os.FileInfo
	)
	// 1. get normal and error  logger of file info
	normalInfo, err = public.NormalWrite.Stat()
	errnoInfo, err = public.ErrnoWrite.Stat()
	if err != nil {
		log.OnWrite = false
		log.Destroy()
		log.Serious(err)
		return false
	}
	// 2. compare size
	errno := public.Round(float64(errnoInfo.Size())/math.Pow(1024, 2)) >= log.MaxSize
	normal := public.Round(float64(normalInfo.Size())/math.Pow(1024, 2)) >= log.MaxSize
	//fmt.Println(float64(errnoInfo.Size())  / math.Pow(1024,2),log.MaxSize)
	if normal || errno {
		if normal {
			atomic.AddUint32(&public.NormalFileCount, 1)
			newSuffix := "." + strconv.Itoa(int(public.NormalFileCount)) + ".log"
			oldStr := public.GlobalRe.FindString(log.LoggerNormal)
			if oldStr == "" {
				oldStr = ".log"
			}
			public.NormalWrite.Close()
			log.LoggerNormal = strings.Replace(log.LoggerNormal, oldStr, newSuffix, -1)
			public.NormalWrite, err = os.OpenFile(log.LoggerNormal, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		} else {
			atomic.AddUint32(&public.ErrorFileCount, 1)
			newSuffix := "." + strconv.Itoa(int(public.ErrorFileCount)) + ".log"
			oldStr := public.GlobalRe.FindString(log.LoggerError)
			if oldStr == "" {
				oldStr = ".log"
			}
			public.ErrnoWrite.Close()
			log.LoggerError = strings.Replace(log.LoggerError, oldStr, newSuffix, -1)
			public.ErrnoWrite, err = os.OpenFile(log.LoggerError, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		}

		if err != nil {
			log.OnWrite = false
			log.Destroy()
			log.Serious(err)
			return false
		}
		// inspect again
		log.FileSizeInspector()
	}
	return true
}

func (log *logger) ShouldBind(PathOrYaml interface{}) (l *logger, err error) {
	var yamlData []byte
	switch PathOrYaml.(type) {
	case string:
		yamlFilePath := PathOrYaml.(string)
		if yamlData, err = ioutil.ReadFile(yamlFilePath); err == nil {
			if err = yaml.Unmarshal(yamlData, log); err == nil {
				return log, nil
			}
		}
		return
	case []byte:
		yamlData := PathOrYaml.([]byte)
		if err = yaml.Unmarshal(yamlData, log); err == nil {
			return log, nil
		}
		return
	default:
		return nil, errors.New("UnKnown Data of Type")
	}
}

func (log *logger) inspector(level string, f interface{}, a ...interface{}) {
	log.lock.Lock()
	a = log.handleOption(a...)
	if !log.OnConsole && !log.OnWrite{
		return
	}
	if public.GlobalLevelInt[log.Level] >= public.GlobalLevelInt[level] {
		log.Msg = public.Format(f, a...)
		log.NowTime = fmt.Sprintf("%v [%s]", time.Now().Format(log.TimeFormat), level)
		log.FilePath = public.FilePath() // current print log the file path and line
		if genTraceID {
			log.TraceID = public.GenTraceID()
		}
		if log.OnConsole {
			log.printRow(level)
		}
		log.writInspector(level)
		log.recoverOption()
	}

}

func (log *logger) printRow(level string) {
	if log.OnColor {
		_, _ = os.Stdout.Write(append([]byte(log.colorBrush(level)), '\n'))
		return
	}
	_, _ = os.Stdout.Write(append([]byte(log.Text()), '\n'))
}

func (log *logger) colorBrush(level string) string {
	// 1. receiver color brush method
	brushColor := public.Colors[public.GlobalLevelInt[level]]
	// 2. change text with color
	nowTime := brushColor(log.NowTime)
	filePath := brushColor(log.FilePath)
	msg := brushColor(log.Identifier + ": " + log.Msg)
	// 3. return has color text
	return fmt.Sprintf("%v %v %v", nowTime, filePath, msg)
}

func (log *logger) Text() string {
	return fmt.Sprintf("%v %v %v", log.NowTime, log.FilePath, log.Identifier+": "+log.Msg)
}

func (log *logger) handleOption(o ...interface{}) (data []interface{}) {
	data = make([]interface{}, 0)
	for _, v := range o {
		switch v.(type) {
		case LoggerHandler:
			if option, ok := v.(LoggerHandler); ok {
				option(log)
			}
		case nil:
			// nothing todo
		default:
			data = append(data, v)
		}
	}

	return
}


func (log *logger) recoverOption() {
	if len(log.option) == 0 {
		return
	}
	for _, v := range log.option {
		switch v {
		case public.CallOffOnWrite:
			log.OnWrite = true
		case public.CallOffOnColor:
			log.OnColor = true
		case public.CallOffOnConsole:
			log.OnConsole = true
		}
	}
	log.option = log.option[:0]
}

func (log *logger) Destroy() {
	public.NormalWrite.Close()
	public.ErrnoWrite.Close()
}
