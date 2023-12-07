package log

import (
	"encoding/json"
	"sync"

	"github.com/zishang520/socket.io/v2/socket"
)

const (
	logPath     = "/var/log/keepakonf.log"
	logPageSize = 10
)

var (
	mu                sync.Mutex
	records           []json.RawMessage
	loadedAllFromFile bool
	io                socket.NamespaceInterface
)

func SetIO(newIO socket.NamespaceInterface) {
	io = newIO
}
