package model

import (
	"errors"
	"fmt"
	"github.com/eliot-jay/logger/public"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type LoggerOutPut interface {
	Debug(f interface{}, a ...interface{})
	Info(f interface{}, a ...interface{})
	Warn(f interface{}, a ...interface{})
	Error(f interface{}, a ...interface{})
	Serious(f interface{}, a ...interface{})
	Fatal(f interface{}, a ...interface{})
	Sprint(f interface{}, a ...interface{}) *Logger
	Destroy()
}




type LoggerConfigure struct {
	OnConsole    bool   `yaml:"on-console"`
	OnColor      bool   `yaml:"on-color"`
	OnWrite      bool   `yaml:"on-write"`
	LoggerNormal string `yaml:"normal-file"`
	LoggerError  string `yaml:"error-file"`
	Level        string `yaml:"level"`
	Identifier   string `yaml:"identifier"`
	TimeFormat   string `yaml:"time-format"`
	MaxSize      int64  `yaml:"max-size"`
}

type Logger struct {
	lock            *sync.Mutex
	TraceID         string
	FilePath        string
	NowTime         string
	Msg             string
	LoggerConfigure `yaml:"loggerConfigure"`
}

func NewDefaultLogger() (log *Logger) {
	log = &Logger{lock: &sync.Mutex{}}
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

func (log *Logger) Debug(f interface{}, a ...interface{}) {
	log.levelInspector(public.DEBUG, f, a...)
}

func (log *Logger) Info(f interface{}, a ...interface{}) {
	log.levelInspector(public.INFO, f, a...)
}

func (log *Logger) Warn(f interface{}, a ...interface{}) {
	log.levelInspector(public.WARN, f, a...)
}

func (log *Logger) Error(f interface{}, a ...interface{}) {
	log.levelInspector(public.ERRNO, f, a...)
}

func (log *Logger) Serious(f interface{}, a ...interface{}) {
	log.levelInspector(public.SERIOUS, f, a...)
}
func (log *Logger) Fatal(f interface{}, a ...interface{}) {
	defer os.Exit(0)
	log.levelInspector(public.FATAL, f, a...)
}

func (log *Logger) Sprint(f interface{}, a ...interface{}) *Logger{
	log.OnConsole = false
	defer func() {
		log.OnConsole = true
	}()
	log.levelInspector(public.Sprint, f, a...)
	return log
}



func (log *Logger) writInspector(level string) {
	if log.OnWrite {
		if log.openFile() {
			log.write(level)
		}

	}
}

func (log *Logger) openFile() bool {
	if public.Once {
		public.Once = false
		var err error
		public.NormalWrite, err = os.OpenFile(log.LoggerNormal, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		public.ErrnoWrite, err = os.OpenFile(log.LoggerError, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal("Open logger File failure: %v", err)
			return false
		}

	}
	return true
}

func (log *Logger) write(level string) {
		log.FileSizeInspector()
		if level == public.ERRNO || level == public.SERIOUS || level == public.FATAL{
			public.ErrnoWrite.WriteString(fmt.Sprintf("%v TraceID: %v \n",log.Text(),log.TraceID))
		} else {
			public.NormalWrite.WriteString(fmt.Sprintf("%v TraceID: %v \n",log.Text(),log.TraceID))
		}
}

func (log *Logger) FileSizeInspector() {
	normalInfo, _ := public.NormalWrite.Stat()
	errnoInfo, _ := public.ErrnoWrite.Stat()
	normal := (normalInfo.Size()/1024)/1024 >= log.MaxSize
	errno := (errnoInfo.Size()/1024)/1024 >= log.MaxSize
	if normal || errno {
		atomic.AddInt64(&public.FileCount,1)
		newSuffix := "." + strconv.Itoa(int(public.FileCount)) + ".log"
		re := regexp.MustCompile("\\.[0-9]+.*?log")
		if normal {
			oldStr := re.FindString(log.LoggerNormal)
			if oldStr == "" {
				oldStr = ".log"
			}
			public.NormalWrite.Close()
			log.LoggerNormal = strings.Replace(log.LoggerNormal, oldStr, newSuffix, -1)
			public.NormalWrite, _ = os.OpenFile(log.LoggerNormal, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		} else {
			oldStr := re.FindString(log.LoggerError)
			if oldStr == "" {
				oldStr = ".log"
			}
			public.ErrnoWrite.Close()
			log.LoggerError = strings.Replace(log.LoggerError, oldStr, newSuffix, -1)
			public.ErrnoWrite, _ = os.OpenFile(log.LoggerError, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		}
		log.FileSizeInspector()
	}

}

func (log *Logger) ShouldBind(PathOrYaml interface{}) (l *Logger, err error) {
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

func (log *Logger) levelInspector(level string, f interface{}, a ...interface{}) {
	log.lock.Lock()
	if public.GlobalLevelInt[log.Level] >= public.GlobalLevelInt[level] {
		log.Msg = public.Format(f, a...)
		log.NowTime = fmt.Sprintf("%v [%s]", time.Now().Format(log.TimeFormat), level)
		log.FilePath = public.FilePath()
		log.TraceID = public.GenTraceID()
		log.writInspector(level)
		if log.OnConsole {
			log.printRow(level)
		}
	}
	log.lock.Unlock()
}

func (log *Logger) printRow(level string) {
	if log.OnColor {
		os.Stdout.Write(append([]byte(log.colorBrush(level)), '\n'))
		return
	}
	os.Stdout.Write(append([]byte(log.Text()),'\n'))
}

func (log *Logger) colorBrush(level string) string {
	// 1. 获取染色函数
	brushColor := public.Colors[public.GlobalLevelInt[level]]
	// 2. 对需要的字符串进行染色
	nowTime := brushColor(log.NowTime)
	filePath := brushColor(log.FilePath)
	msg := brushColor(log.Identifier + ": " + log.Msg)
	// 3. 返回已经染色的字符串
	return fmt.Sprintf("%v %v %v", nowTime, filePath, msg)
}

func (log *Logger) Text() string {
	return fmt.Sprintf("%v %v %v", log.NowTime, log.FilePath, log.Identifier+": "+log.Msg)
}

func (log *Logger)Destroy()  {
	public.NormalWrite.Close()
	public.ErrnoWrite.Close()
}