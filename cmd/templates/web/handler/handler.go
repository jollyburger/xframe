package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

var Router = mux.NewRouter()

type HandleFunc struct {
	Handler func(http.ResponseWriter, *http.Request)
	Methods []string
}

var handlerMap = map[string]HandleFunc{
	"/": HandleFunc{
		Handler: hello,
		Methods: []string{"GET"},
	},
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

func RegisterHandler() *mux.Router {
	for url, handler := range handlerMap {
		Router.HandleFunc(url, handler.Handler).Methods(handler.Methods...)
	}
	return Router
}
