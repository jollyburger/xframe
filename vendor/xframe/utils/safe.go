package utils

import (
	"net/http"
	"runtime"
	"xframe/log"
	"xframe/trace"
)

var (
	panicProtection bool = false
)

func InitWrapPanic(flag bool) {
	panicProtection = flag
}

func CheckWrapPanic() bool {
	return panicProtection
}

func GoSafeTCP(ctx trace.XContext, ch chan []byte, req interface{}, fn func(trace.XContext, chan []byte, interface{})) {
	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, 1024*8)
			stack = stack[:runtime.Stack(stack, false)]
			f := "[PANIC] %s\n%s"
			log.ERRORF(f, err, stack)
		}
	}()
	fn(ctx, ch, req)
}

func GoSafeHTTP(ctx trace.XContext, rw http.ResponseWriter, r *http.Request, fn func(trace.XContext, http.ResponseWriter, *http.Request)) {
	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, 1024*8)
			stack = stack[:runtime.Stack(stack, false)]
			f := "[PANIC] %s\n%s"
			log.ERRORF(f, err, stack)
		}
	}()
	fn(ctx, CtxInHttpRspHeader(ctx, rw), CtxInHttpReqHeader(ctx, r))
}

func GoSafeTimer(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, 1024*8)
			stack = stack[:runtime.Stack(stack, false)]
			f := "[PANIC] %s\n%s"
			log.ERRORF(f, err, stack)
		}
	}()
	fn()
}

func GoSafe(fn func(...interface{}), args ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, 1024*8)
			stack = stack[:runtime.Stack(stack, false)]
			f := "[PANIC] %s\n%s"
			log.ERRORF(f, err, stack)
		}
	}()
	fn(args...)
}
