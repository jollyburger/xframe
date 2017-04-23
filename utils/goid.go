package utils

import (
	"bytes"
	"runtime"
	"strconv"
)

func GetGID() (gid uint64) {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	gid, _ = strconv.ParseUint(string(b), 10, 64)
	return
}
