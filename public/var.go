package public

import "os"

var (
	GlobalLevelInt  = make(map[string]int,10)
	NormalWrite *os.File
	ErrnoWrite *os.File
	Once = true
	FileCount int64
)