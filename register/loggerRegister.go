package register

var genTraceID = true


type LoggerRegister struct {
	log *logger
}

func NewDefaultLogger() Logger {
	log := &LoggerRegister{}
	return log.defaultLoggerRegister()
}

func NewLoggerBy(yamlPath string) (Logger, error) {
	log := &LoggerRegister{}
	return log.newLoggerRegister(yamlPath)
}

func (r *LoggerRegister) Debug(f interface{}, a ...interface{}) string {
	return r.log.Debug(f, a...)
}

func (r *LoggerRegister) Info(f interface{}, a ...interface{}) string{
	return r.log.Info(f, a...)
}

func (r *LoggerRegister) Warn(f interface{}, a ...interface{}) string{
	return r.log.Warn(f, a...)
}

func (r *LoggerRegister) Error(f interface{}, a ...interface{}) string{
	return r.log.Error(f, a...)
}

func (r *LoggerRegister) Serious(f interface{}, a ...interface{})string {
	return r.log.Serious(f, a...)
}

func (r *LoggerRegister) Fatal(f interface{}, a ...interface{}) {
	r.log.Fatal(f, a...)
}

// self define type field of highest level
// the type print of non-color
func (r *LoggerRegister) sPrintf(TractID,Type string, f interface{}, a ...interface{}) string {
	return r.log.sPrintf(TractID,Type, f, a...)
}


func (r *LoggerRegister) defaultLoggerRegister() Logger {
    return r.init()
}

func (r *LoggerRegister) newLoggerRegister(yamlPath string) (l Logger, err error) {
	r.log, err = r.init().log.ShouldBind(yamlPath)
	return r, err
}

func (r *LoggerRegister) init() *LoggerRegister {
	r.log = newDefaultLogger()
	return r
}
func (r *LoggerRegister) Destroy() {
	r.log.Destroy()
}
