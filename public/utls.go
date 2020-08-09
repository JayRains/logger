package public

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type brush func(string) string

var (
	Colors = []brush{
		newBrush("1;41"), // Serious              红色底	0
		newBrush("1;31"), // Error                红色		1
		newBrush("1;35"), // Warn               紫红底	2
		newBrush("1;34"), // Info               蓝色		3
		newBrush("1;32"), // Debug              绿色		4
		newBrush("4;36"), //underline          青色+下划线 5
	}
)

func init() {
	//log level
	GlobalLevelInt[DEBUG] = 4
	GlobalLevelInt[INFO] = 3
	GlobalLevelInt[WARN] = 2
	GlobalLevelInt[ERRNO] = 1
	GlobalLevelInt[SERIOUS] = 0
	GlobalLevelInt[FATAL] = 0
	GlobalLevelInt[Underline] = 5
}

func newBrush(color string) brush {
	prefix := "\033[" // \033[ 1; 32m%s  \033[0m
	suffix := "\033[0m"
	return func(text string) string {
		return prefix + color + "m" + text + suffix
	}

}

func Format(f interface{}, v ...interface{}) string {
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

func FilePath() string {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(Skip, rpc[:])
	if n < 1 {
		return ""
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	currentDir, _ := os.Getwd()
	return strings.Replace(frame.File, currentDir+"/", "", -1) + ToString(frame.Line)
}

func GenTraceID() string {
	b := bytes.Buffer{}
	// 纳秒时间
	nano := time.Now().UnixNano()
	// 纳秒随机数种子时间
	rand.Seed(nano)
	randomNum := rand.Int()
	// 不断变化的pid
	pid := os.Getpid()
	// cpu 的使用率
	percent, _ := cpu.Percent(time.Nanosecond, false)
	// 内存 的使用率
	v, _ := mem.VirtualMemory()
	MemoryUsedPercent := fmt.Sprint(v.UsedPercent)
	CpuUsedPercent := fmt.Sprint(percent)
	// 基础数据
	base := strconv.Itoa(int(nano)) + strconv.Itoa(randomNum) + strconv.Itoa(pid)
	// 合成所有因素
	base += MemoryUsedPercent + CpuUsedPercent
	hash := sha1.New()
	hash.Write([]byte(base))
	b.WriteString(fmt.Sprintf("%x", hash.Sum(nil)))
	return b.String()
}
