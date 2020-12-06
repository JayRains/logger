package option

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/gookit/color"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"math"
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
		NewBrush(Red, true), // FATAL SERIOUS 红色底        0
		NewBrush(Red),       // Error              红色			1
		NewBrush(Purple),    // Warn               紫红		    2
		NewBrush(Blue),      // Info               蓝色			3
		NewBrush(Green),     // Debug              绿色			4
	}
)

func init() {
	//log level
	globallevelint[SERIOUS] = 0
	globallevelint[FATAL] = 0
	globallevelint[ERRNO] = 1
	globallevelint[WARN] = 2
	globallevelint[INFO] = 3
	globallevelint[DEBUG] = 4

}

func NewBrush(hex string, isBG ...bool) brush {
	return func(text string) string {
		return color.HEX(hex, isBG...).Sprint(text)
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

func FilePath() string {
	skip := Skip
	rpc := make([]uintptr, 1)
	if register {
		skip++
	}
	n := runtime.Callers(skip, rpc[:])
	if n < 1 {
		return ""
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	currentDir, _ := os.Getwd()
	return strings.Replace(frame.File, currentDir+"/", "", -1) + ":" + strconv.Itoa(frame.Line)
}

func GenTraceID() string {
	b := bytes.Buffer{}
	// 纳秒时间
	nano := time.Now().UnixNano()
	// 纳秒随机数种子时间
	rand.Seed(nano)
	randomNum := rand.Int()
	// 变化的pid
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
	b.WriteString(fmt.Sprintf("%x ", hash.Sum(nil)))
	return b.String()
}

// 0.245 => 0.25
func Round(x float64) float64 {
	n10 := math.Pow10(2)
	return math.Trunc((x+0.5/n10)*n10) / n10
}
