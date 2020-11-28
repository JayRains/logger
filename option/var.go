package option

import (
	"os"
	"regexp"
	"sync"
)

var (
	genTraceID = true
	register   = false
	globallevelint  = make(map[string]int, 10)
	normalwrite     *os.File
	errnowrite      *os.File
	once            = &sync.Once{}
	normalfilecount uint32
	errorfilecount  uint32
	globalre        = regexp.MustCompile("\\.[0-9]+.*?log")
)
