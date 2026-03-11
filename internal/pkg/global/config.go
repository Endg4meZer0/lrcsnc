package global

import (
	"sync"

	configStruct "lrcsnc/internal/config/struct"
)

var Config = struct {
	M sync.Mutex
	C configStruct.Config

	Path string
}{}

// Version is linked through -X (check Makefile)
var Version = "dev"
