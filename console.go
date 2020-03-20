package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const (
	//level field
	debug   = "DEBUG"
	info    = "INFO"
	warn    = "WARN"
	err     = "ERROR"
	serious = "SERIOUS"


	//system variable
	skip      = 5
	underline = "Underline"
	logChannelCache = 100
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
	levelInt = make(map[string]int)
	log      = make(chan string, logChannelCache)
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
