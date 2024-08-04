package ws

import (
	"os"
	"time"
)

func ReadFileIfModified(filename string, lastMod time.Time) ([]byte, time.Time, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, lastMod, err
	}

	modTime := fi.ModTime()
	if !modTime.After(lastMod) {
		return nil, lastMod, nil
	}

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, lastMod, err
	}
	return bytes, modTime, nil
}
