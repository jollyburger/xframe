package utils

import (
	"fmt"
	"net/http"
	"xframe/trace"
)

const (
	CONTEXT_KEY = "X-Trace-Context"
)

func CtxInHttpReqHeader(ctx trace.XContext, r *http.Request) *http.Request {
	r.Header.Set(CONTEXT_KEY, fmt.Sprintf("%d:%d:%s", ctx.GetTraceId(), ctx.GetSpanId(), ctx.GetSessionNo()))
	return r
}

func CtxInHttpRspHeader(ctx trace.XContext, rw http.ResponseWriter) http.ResponseWriter {
	rw.Header().Set(CONTEXT_KEY, fmt.Sprintf("%d", ctx.GetSpanId()))
	return rw
}
