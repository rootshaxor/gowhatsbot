package utils

import (
	"path"
	"runtime"
)

func GetMyPath() string {
	if _, filename, _, ok := runtime.Caller(1); ok {
		return path.Dir(filename)
	}

	return ""
}
