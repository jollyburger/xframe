package utils

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"xframe/trace"
)

const (
	CONTEXT_KEY = "X-Trace-Context"
)

func CtxInHttpReqHeader(ctx trace.XContext, r *http.Request) *http.Request {
	//r.Header.Set(CONTEXT_KEY, fmt.Sprintf("%d:%d:%s", ctx.GetTraceId(), ctx.GetSpanId(), ctx.GetSessionNo()))
	return r
}

func CtxInHttpRspHeader(ctx trace.XContext, rw http.ResponseWriter) http.ResponseWriter {
	//rw.Header().Set(CONTEXT_KEY, fmt.Sprintf("%d", ctx.GetSpanId()))
	return rw
}

//v0.3
func convertType(sf reflect.StructField, rv reflect.Value, value string) {
	var v interface{}
	switch sf.Type.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32:
		tmp_v, err := strconv.Atoi(value)
		if err != nil {
			return
		}
		v = tmp_v
	case reflect.Int64:
		tmp_v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return
		}
		v = tmp_v
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		tmp_v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return
		}
		v = tmp_v
	case reflect.Float32, reflect.Float64:
		tmp_v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return
		}
		v = tmp_v
	case reflect.String:
		v = value
	case reflect.Bool:
		v = (value == "1")
	}
	rv.Set(reflect.ValueOf(v))
}

func ParseGetRequest(req *http.Request, obj interface{}, tag string) {
	v := reflect.Indirect(reflect.ValueOf(obj))
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		rv := v.Field(i)
		if _, ok := sf.Tag.Lookup(tag); !ok {
			continue
		}
		st := sf.Tag.Get(tag)
		reqV := req.FormValue(st)
		if reqV != "" {
			convertType(sf, rv, reqV)
		}
	}
}

//json body
func ParsePostRequest(req *http.Request, obj interface{}) (err error) {
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&obj)
	if err != nil {
		return
	}
	defer req.Body.Close()
	return
}
