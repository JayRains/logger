package register

import (
	"errors"
	"fmt"
	"github.com/eliot-jay/logger/public"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
	sPrintf(TrackID,Type string, f interface{}, a ...interface{}) string
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
	lock            *sync.RWMutex
	FilePath        string
	TraceID    		string
	NowTime         string
	Msg             string
	LoggerConfigure `yaml:"loggerConfigure"`
}

func newDefaultLogger() (log *logger) {
	log = &logger{lock: &sync.RWMutex{}}
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

func (log *logger) Debug(f interface{}, a ...interface{})string{
	defer log.lock.Unlock()
	log.levelInspector(public.DEBUG, f, a...)
	return log.TraceID
}

func (log *logger) Info(f interface{}, a ...interface{})string{
	defer log.lock.Unlock()
	log.levelInspector(public.INFO, f, a...)
	return log.TraceID
}

func (log *logger) Warn(f interface{}, a ...interface{})string{
	defer log.lock.Unlock()
	log.levelInspector(public.WARN, f, a...)
	return log.TraceID
}

func (log *logger) Error(f interface{}, a ...interface{})string{
	defer log.lock.Unlock()
	log.levelInspector(public.ERRNO, f, a...)
	return log.TraceID
}

func (log *logger) Serious(f interface{}, a ...interface{})string {
	defer log.lock.Unlock()
	log.levelInspector(public.SERIOUS, f, a...)
	return log.TraceID
}
func (log *logger) Fatal(f interface{}, a ...interface{})  {
	defer os.Exit(0)
	defer log.lock.Unlock()
	log.levelInspector(public.FATAL, f, a...)
}

func (log *logger) sPrintf(TractID,Type string, f interface{}, a ...interface{}) string {
	defer log.lock.Unlock()
	genTraceID = false
	log.TraceID = TractID
	if log.OnConsole {
		log.OnConsole = false
		defer func() {
			log.OnConsole = true
		}()
	}
	log.levelInspector(Type, f, a...)
	genTraceID = true
	return log.Text()
}

func (log *logger) writInspector(level string) {
	log.FileSizeInspector()
	if log.OnWrite {
		log.write(level)
	}
}

func (log *logger) write(level string) {
	if level == public.ERRNO || level == public.SERIOUS || level == public.FATAL {
		public.ErrnoWrite.WriteString(fmt.Sprintf("%v %v: %v \n", log.Text(),public.Trace,log.TraceID))
	} else {
		public.NormalWrite.WriteString(fmt.Sprintf("%v %v: %v \n", log.Text(),public.Trace,log.TraceID))
	}
}

// logger file size compare
func (log *logger) FileSizeInspector() {
	if !log.OnWrite {
		return
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
	// 2. compare size
	normal := (float64(normalInfo.Size()/1024) / 1024) >= log.MaxSize
	errno := (float64(errnoInfo.Size()/1024) / 1024) >= log.MaxSize
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
			return
		}
		// inspect again
		log.FileSizeInspector()
	}

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

func (log *logger) levelInspector(level string, f interface{}, a ...interface{}) {
	log.lock.Lock()
	if public.GlobalLevelInt[log.Level] >= public.GlobalLevelInt[level] {
		log.Msg = public.Format(f, a...)
		log.NowTime = fmt.Sprintf("%v [%s]", time.Now().Format(log.TimeFormat), level)
		log.FilePath = public.FilePath() // current print log the file path
		if genTraceID {  log.TraceID = public.GenTraceID()  }
		log.writInspector(level)
		if log.OnConsole {
			log.printRow(level)
		}
	}

}

func (log *logger) printRow(level string) {
	if log.OnColor {
		os.Stdout.Write(append([]byte(log.colorBrush(level)), '\n'))
		return
	}
	os.Stdout.Write(append([]byte(log.Text()), '\n'))
}

func (log *logger) colorBrush(level string) string {
	// 1. 获取染色函数
	brushColor := public.Colors[public.GlobalLevelInt[level]]
	// 2. 对需要的字符串进行染色
	nowTime := brushColor(log.NowTime)
	filePath := brushColor(log.FilePath)
	msg := brushColor(log.Identifier + ": " + log.Msg)
	// 3. 返回已经染色的字符串
	return fmt.Sprintf("%v %v %v", nowTime, filePath, msg)
}

func (log *logger) Text() string {
	return fmt.Sprintf("%v %v %v", log.NowTime, log.FilePath, log.Identifier+": "+log.Msg)
}

func (log *logger) Destroy() {
	public.NormalWrite.Close()
	public.ErrnoWrite.Close()
}
