package public

import (
	"os"
	"regexp"
	"sync"
)

var (
	GlobalLevelInt  = make(map[string]int, 10)
	NormalWrite     *os.File
	ErrnoWrite      *os.File
	Once            = &sync.Once{}
	NormalFileCount uint32
	ErrorFileCount  uint32
	GlobalRe        = regexp.MustCompile("\\.[0-9]+.*?log")
)
