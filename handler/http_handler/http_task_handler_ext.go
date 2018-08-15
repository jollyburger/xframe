package http_handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

//TODO: customize code
var (
	SUCCESS = 0
)

//TODO: add trace
var (
	Rt             = mux.NewRouter()
	DoBaseResponse = func(http.ResponseWriter, int) {}
	DoDataResponse = func(http.ResponseWriter, int, interface{}) {}
)

func httpWrapper(f func(*http.Request) (interface{}, int)) func(http.ResponseWriter, *http.Request) {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		data, err := f(r)
		if err != SUCCESS {
			DoBaseResponse(rw, err)
		} else {
			if data != nil {
				DoDataResponse(rw, SUCCESS, data)
			} else {
				DoBaseResponse(rw, SUCCESS)
			}
		}
	}
	return fn
}

func RegisterHTTPMuxHandler(path string, f func(*http.Request) (interface{}, int)) *mux.Route {
	return Rt.HandleFunc(path, httpWrapper(f))
}
