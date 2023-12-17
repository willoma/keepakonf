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

func GetPage(offset int) logMessage {
	mu.Lock()
	defer mu.Unlock()

	needToHave := offset + logPageSize

	if len(records) < needToHave {
		// We do not have enough records...

		// If we have not already loaded the whole log file, try to read more.
		if !loadedAllFromFile {
			loadFromFile(needToHave)
		}

		// If we have read everything from the file, then do not return more
		// than what we have.
		if loadedAllFromFile {
			end := len(records) - offset
			if end < 0 {
				return logMessage{[]json.RawMessage{}, true}
			}
			return logMessage{records[:end], true}
		}
	}

	// We have enough records to provide the requested page
	start := len(records) - needToHave

	return logMessage{records[start : start+logPageSize], false}
}

func loadFromFile(count int) {
	f, err := os.Open(logPath)
	if err != nil {
		slog.Error("Could not open log file", "error", err)
		return
	}
	defer f.Close()

	fs, err := f.Stat()
	if err != nil {
		slog.Error("Could not get stats on log file", "error", err)
		return
	}

	newRecords := make([]json.RawMessage, 0, count)

	sc := rscanner.NewScanner(f, fs.Size())
	i := 0
	for i < count && sc.Scan() {
		size := len(sc.Bytes())
		if size == 0 {
			continue
		}
		dst := make([]byte, size)
		copy(dst, sc.Bytes())
		newRecords = append(newRecords, json.RawMessage(dst))
		i++
	}

	if sc.Err() != nil {
		slog.Error("Could not read logs from file", "error", sc.Err())
		return
	}

	if len(newRecords) < count {
		loadedAllFromFile = true
	}

	slices.Reverse(newRecords)

	records = newRecords
}
