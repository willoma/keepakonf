package data

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/zishang520/socket.io/v2/socket"

	"github.com/willoma/keepakonf/internal/log"
	"github.com/willoma/keepakonf/internal/runners"
)

const (
	dbdirpath  = "/var/lib/keepakonf"
	dbfilepath = dbdirpath + "/keepakonf.db"
	dedupDelay = 500 * time.Millisecond
)

type Data struct {
	groups []*runners.Group
	mu     sync.Mutex

	io *socket.Server

	dedupTimer *time.Timer
	dedupMutex sync.Mutex
}

func New(io *socket.Server) *Data {
	d := &Data{
		io: io,
	}
	d.load()
	return d
}

func (d *Data) load() {
	f, err := os.ReadFile(dbfilepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			d.groups = []*runners.Group{}
			log.Info(
				"No database file, starting from scratch",
				"database",
				"", "", "", "", nil,
			)
			return
		}
		log.Errorf(err, "Could not read database file %s", dbfilepath)
		return
	}

	var src []any
	if err := json.Unmarshal(f, &src); err != nil {
		log.Errorf(err, "Wrong data in the database file %s", dbfilepath)
		return
	}

	d.mu.Lock()
	d.groups = make([]*runners.Group, len(src))
	for i, srcG := range src {
		g := runners.GroupFromMap(srcG, d.io.Sockets())
		g.Watch()
		d.groups[i] = g
	}
	d.mu.Unlock()
}

func (d *Data) save() {
	d.dedupMutex.Lock()
	if d.dedupTimer == nil {
		d.dedupTimer = time.AfterFunc(dedupDelay, func() {
			d.doSave()
			d.dedupTimer = nil
		})
	} else {
		d.dedupTimer.Reset(dedupDelay)
	}
	d.dedupMutex.Unlock()
}

func (d *Data) doSave() {
	d.mu.Lock()
	dst := make([]any, len(d.groups))
	for i, grp := range d.groups {
		instructionsClone := make([]any, len(grp.Instructions))
		for j, ins := range grp.Instructions {
			paramsClone := make(map[string]any, len(ins.Parameters))
			for k, v := range ins.Parameters {
				paramsClone[k] = v
			}
			instructionsClone[j] = map[string]any{
				"id":         ins.ID,
				"command":    ins.Command,
				"parameters": paramsClone,
			}
		}
		dst[i] = map[string]any{
			"id":           grp.ID,
			"name":         grp.Name,
			"instructions": instructionsClone,
		}
	}
	d.mu.Unlock()

	bin, err := json.Marshal(dst)
	if err != nil {
		log.Error(err, "Wrong data before saving the configuration")
		return
	}

	if err := os.MkdirAll(dbdirpath, 0700); err != nil {
		log.Error(err, "Could not create configuration directory")
		return
	}
	if err := os.WriteFile(dbfilepath, bin, 0644); err != nil {
		log.Error(err, "Could not save configuration")
	}
}
