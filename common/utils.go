package common

import (
	"path/filepath"
	"runtime"
)

func GetCurDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}
