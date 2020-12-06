package option

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Configure struct {
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


type Logger struct {
	FilePath  string
	TraceID   string
	NowTime   string
	Msg       string
	Configure `yaml:"loggerConfigure"`

	//private
	lock   *sync.Mutex
	option []uint
}

func Default() (log *Logger) {
	log = &Logger{lock: &sync.Mutex{}, option: make([]uint, 0)}
	config := `
loggerConfigure:
  on-color: true  
  on-write: false 
  on-console: true
  level: "DEBUG"
  max-size: 1
  normal-file: "./logs/project.inf.1.log"
  error-file: "./logs/project.err.1.log"
  identifier: "" 
  time-format: "2006-01-02 15:04:05"
`
	log, _ = log.ShouldBind([]byte(config))
	return log
}

func NewLoggerBy(yamlPath string) (*Logger, error) {
	log := Default()
	return log.ShouldBind(yamlPath)
}

func (log *Logger) Debug(f interface{}, a ...interface{}) string {
	log.inspector(DEBUG, f, a...)
	return log.TraceID
}

func (log *Logger) Info(f interface{}, a ...interface{}) string {
	log.inspector(INFO, f, a...)
	return log.TraceID
}

func (log *Logger) Warn(f interface{}, a ...interface{}) string {
	log.inspector(WARN, f, a...)
	return log.TraceID
}

func (log *Logger) Error(f interface{}, a ...interface{}) string {
	log.inspector(ERRNO, f, a...)
	return log.TraceID
}

func (log *Logger) Serious(f interface{}, a ...interface{}) string {
	log.inspector(SERIOUS, f, a...)
	return log.TraceID
}

func (log *Logger) Fatal(f interface{}, a ...interface{}) {
	defer os.Exit(0)
	defer log.Destroy()
	log.inspector(FATAL, f, a...)
}

func (log *Logger) SPrintf(TractID, Type string, f interface{}, a ...interface{}) string {
	log.lock.Lock()
	if log.OnConsole {
		log.OnConsole = false
		defer func() {
			log.OnConsole = true
		}()
	}
	genTraceID = false
	log.TraceID = TractID
	log.lock.Unlock()
	log.inspector(Type, f, a...)
	genTraceID = true

	return log.Text()
}

func (log *Logger) writInspector(level string) {
	if log.FileSizeInspector() {
		log.write(level)
	}
}

func (log *Logger) write(level string) {
	if level == ERRNO || level == SERIOUS || level == FATAL {
		_, _ = errnowrite.WriteString(fmt.Sprintf("%v %v: %v \n", log.Text(), Trace, log.TraceID))
	} else {
		_, _ = normalwrite.WriteString(fmt.Sprintf("%v %v: %v \n", log.Text(), Trace, log.TraceID))
	}
}

// option file size compare
func (log *Logger) FileSizeInspector() bool {
	if !log.OnWrite {
		return false
	}
	once.Do(func() {
		var err error
		normalwrite, err = os.OpenFile(log.LoggerNormal, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		errnowrite, err = os.OpenFile(log.LoggerError, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal("Open option File failure: %v", err)
		}
	})

	var (
		err        error
		normalInfo os.FileInfo
		errnoInfo  os.FileInfo
	)
	// 1. get normal and error  option of file info
	normalInfo, err = normalwrite.Stat()
	errnoInfo, err = errnowrite.Stat()
	if err != nil {
		log.OnWrite = false
		log.Destroy()
		log.Serious(err)
		return false
	}
	// 2. compare size
	errno := Round(float64(errnoInfo.Size())/math.Pow(1024, 2)) >= log.MaxSize
	normal := Round(float64(normalInfo.Size())/math.Pow(1024, 2)) >= log.MaxSize
	//fmt.Println(float64(errnoInfo.Size())  / math.Pow(1024,2),log.MaxSize)
	if normal || errno {
		if normal {
			atomic.AddUint32(&normalfilecount, 1)
			newSuffix := "." + strconv.Itoa(int(normalfilecount)) + ".log"
			oldStr := globalre.FindString(log.LoggerNormal)
			if oldStr == "" {
				oldStr = ".log"
			}
			normalwrite.Close()
			log.LoggerNormal = strings.Replace(log.LoggerNormal, oldStr, newSuffix, -1)
			normalwrite, err = os.OpenFile(log.LoggerNormal, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		} else {
			atomic.AddUint32(&errorfilecount, 1)
			newSuffix := "." + strconv.Itoa(int(errorfilecount)) + ".log"
			oldStr := globalre.FindString(log.LoggerError)
			if oldStr == "" {
				oldStr = ".log"
			}
			errnowrite.Close()
			log.LoggerError = strings.Replace(log.LoggerError, oldStr, newSuffix, -1)
			errnowrite, err = os.OpenFile(log.LoggerError, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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

func (log *Logger) inspector(level string, f interface{}, a ...interface{}) {
	log.lock.Lock()
	defer log.lock.Unlock()
	a = log.handleOption(a...)
	if !log.OnConsole && !log.OnWrite {
		return
	}
	if globallevelint[log.Level] >= globallevelint[level] {
		log.Msg = Format(f, a...)
		log.NowTime = fmt.Sprintf("%v [%s]", time.Now().Format(log.TimeFormat), level)
		log.FilePath = FilePath() // current print log the file path and line
		if genTraceID {
			log.TraceID = GenTraceID()
		}
		if log.OnConsole {
			log.printRow(level)
		}
		log.writInspector(level)
		log.recoverOption()
	}

}

func (log *Logger) printRow(level string) {
	if log.OnColor {
		_, _ = os.Stdout.Write(append([]byte(log.colorBrush(level)), '\n'))
		return
	}
	_, _ = os.Stdout.Write(append([]byte(log.Text()), '\n'))
}

func (log *Logger) colorBrush(level string) string {
	// 1. receiver color brush method
	brushColor := Colors[globallevelint[level]]
	// 2. change text with color
	nowTime := brushColor(log.NowTime)
	filePath := brushColor(log.FilePath)
	msg := brushColor(log.Identifier + ": " + log.Msg)
	// 3. return has color text
	return fmt.Sprintf("%v %v %v", nowTime, filePath, msg)
}

func (log *Logger) Text() string {
	return fmt.Sprintf("%v %v %v", log.NowTime, log.FilePath, log.Identifier+": "+log.Msg)
}

func (log *Logger) handleOption(o ...interface{}) (data []interface{}) {
	data = make([]interface{}, 0)
	for _, v := range o {
		switch v.(type) {
		case LogOption:
			if option, ok := v.(LogOption); ok {
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

func (log *Logger) recoverOption() {
	if len(log.option) == 0 {
		return
	}
	for _, v := range log.option {
		switch v {
		case CallOffOnWrite:
			log.OnWrite = true
		case CallOffOnColor:
			log.OnColor = true
		case CallOffOnConsole:
			log.OnConsole = true
		}
	}
	log.option = log.option[:0]
}

func (log *Logger) Destroy() {
	normalwrite.Close()
	errnowrite.Close()
}

func (log *Logger) Register(level, hexColor string, id int) func(f string, v ...interface{}) {
	if id > 4 {
		log.Error("register must id <= 4 ")
		return nil
	}
	log.lock.Lock()
	globallevelint[level] = id
	Colors[id] = NewBrush(hexColor)
	log.lock.Unlock()
	return log.Println(level)
}
func (log *Logger) Println(level string) func(f string, v ...interface{}) {
	register = true
	return func(f string, v ...interface{}) {
		v = append(v, OnConsole())
		log.SPrintf(GenTraceID(), level, f, v...)
	}
}
