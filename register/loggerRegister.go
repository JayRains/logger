package register

import "github.com/eliot-jay/logger/model"

type LoggerRegister struct {
	log *model.Logger
}

func NewDefaultLogger() *LoggerRegister {
	log := &LoggerRegister{}
	return log.defaultLoggerRegister()
}

func NewLogger(yamlPath string) (*LoggerRegister,error)  {
	log := &LoggerRegister{}
	return log.newLoggerRegister(yamlPath)
}


func (r *LoggerRegister) Debug(f interface{}, a ...interface{}) {
	r.log.Debug(f,a...)
}

func (r *LoggerRegister) Info(f interface{}, a ...interface{}) {
	r.log.Info(f,a...)
}

func (r *LoggerRegister) Warn(f interface{}, a ...interface{}) {
	r.log.Warn(f,a...)
}

func (r *LoggerRegister) Error(f interface{}, a ...interface{}) {
	r.log.Error(f,a...)
}

func (r *LoggerRegister) Serious(f interface{}, a ...interface{}) {
	r.log.Serious(f,a...)
}

func (r *LoggerRegister) Fatal(f interface{}, a ...interface{}) {
	r.log.Fatal(f,a...)
}

func (r *LoggerRegister) Sprint(Type string,f interface{}, a ...interface{}) *model.Logger {
	return r.log.Sprint(Type,f,a...)
}




func (r *LoggerRegister)defaultLoggerRegister() *LoggerRegister {
	return r.init()
}

func (r *LoggerRegister)newLoggerRegister(yamlPath string)(l *LoggerRegister  , err error) {
	r.log ,err = r.init().log.ShouldBind(yamlPath)
	return r , err
}

func (r *LoggerRegister)init() *LoggerRegister  {
	r.log = model.NewDefaultLogger()
	return r
}
func (r *LoggerRegister)Destroy() {
	r.log.Destroy()
}
