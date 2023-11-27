package log

import (
	"encoding/json"
	"log/slog"
	"os"
	"slices"

	"github.com/vallerion/rscanner"
)

type logMessage struct {
	Logs          []json.RawMessage `json:"logs"`
	ReachedTheEnd bool              `json:"reached_the_end"`
}

func (l *Logger) GetPage(offset int) logMessage {
	l.mu.Lock()
	defer l.mu.Unlock()

	needToHave := offset + logPageSize

	if len(l.records) < needToHave {
		// We do not have enough records...

		// If we have not already loaded the whole log file, try to read more.
		if !l.loadedAllFromFile {
			l.loadFromFile(needToHave)
		}

		// If we have read everything from the file, then do not return more
		// than what we have.
		if l.loadedAllFromFile {
			end := len(l.records) - offset
			if end < 0 {
				return logMessage{[]json.RawMessage{}, true}
			}
			return logMessage{l.records[:end], true}
		}
	}

	// We have enough records to provide the requested page
	start := len(l.records) - needToHave

	return logMessage{l.records[start : start+logPageSize], false}
}

func (l *Logger) loadFromFile(count int) {
	f, err := os.Open(l.filepath)
	if err != nil {
		l.Error(err, "Could not read log file")
		slog.Error("Could not open log file", "error", err)
		return
	}
	defer f.Close()

	fs, err := f.Stat()
	if err != nil {
		l.Error(err, "Could not read log file")
		slog.Error("Could not get stats on log file", "error", err)
		return
	}

	records := make([]json.RawMessage, 0, count)

	sc := rscanner.NewScanner(f, fs.Size())
	i := 0
	for sc.Scan() && i < count {
		size := len(sc.Bytes())
		if size == 0 {
			continue
		}
		dst := make([]byte, size)
		copy(dst, sc.Bytes())
		records = append(records, json.RawMessage(dst))
		i++
	}

	if sc.Err() != nil {
		l.Error(err, "Could not read log file")
		slog.Error("Could not read logs from file", "error", sc.Err())
		return
	}

	if len(records) < count {
		l.loadedAllFromFile = true
	}

	slices.Reverse(records)
	l.records = records
}
