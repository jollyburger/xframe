package http_handler

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jollyburger/xframe/trace"
)

type XHttpTaskHandler interface {
	ServeXHTTP(ctx trace.XContext, rw http.ResponseWriter, r *http.Request)
}

type XHTTPTaskFunc func(ctx trace.XContext, rw http.ResponseWriter, r *http.Request)

func (t XHTTPTaskFunc) ServeXHTTP(ctx trace.XContext, rw http.ResponseWriter, r *http.Request) {
	t(ctx, rw, r)
}

type httpTaskHandle struct {
	handler XHttpTaskHandler
	timeOut time.Duration
}

var (
	httpHandlePoolMu sync.Mutex
	httpHandlePool   = make(map[string]*httpTaskHandle)
)

func RegisterHTTPTaskHandle(pattern string, handler XHttpTaskHandler, timeOut time.Duration) {
	httpHandlePoolMu.Lock()
	defer httpHandlePoolMu.Unlock()
	newHandle := &httpTaskHandle{
		handler: handler,
		timeOut: timeOut,
	}
	httpHandlePool[pattern] = newHandle
}

func GetHTTPTaskHandle(pattern string) (*httpTaskHandle, error) {
	httpHandlePoolMu.Lock()
	defer httpHandlePoolMu.Unlock()
	if handle, ok := httpHandlePool[pattern]; ok {
		return handle, nil
	} else {
		return nil, errors.New("can't not find  handle")
	}
}

func DumpHTTPTaskHandle() {
	httpHandlePoolMu.Lock()
	defer httpHandlePoolMu.Unlock()
	for k, v := range httpHandlePool {
		fmt.Println(k, v)
	}
}
