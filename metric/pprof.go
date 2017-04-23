package metric

import (
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"
)

func InitPprof(addr string) {
	http.HandleFunc("/goroutine", func(w http.ResponseWriter, r *http.Request) {
		num := strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
		w.Write([]byte(num))
	})
	http.ListenAndServe(addr, nil)
}
