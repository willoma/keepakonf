package log

import (
	"encoding/json"
	"sync"

	"github.com/zishang520/socket.io/v2/socket"
)

const (
	logPath     = "/var/log/keepakonf.log"
	logPageSize = 5
)

type Logger struct {
	filepath          string
	mu                sync.Mutex
	records           []json.RawMessage
	loadedAllFromFile bool

	io socket.NamespaceInterface
}

func NewLogService(io socket.NamespaceInterface) *Logger {
	return &Logger{
		filepath: logPath,
		io:       io,
	}
}
