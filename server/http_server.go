package server

import (
	"net/http"
	"time"
	"xframe/server/websocket"
)

func listenAndServeHTTP(addr string) error {
	http.HandleFunc("/", RouteHTTP)
	http.Handle("/ws", websocket.Handler(RouteWs))
	return http.ListenAndServe(addr, nil)
}

/*
customized multiplexer，Usage:
	mux includes:
	mux.HandleFunc("/", RouteHTTP)
	mux.Handle("/ws", websocket.Handler(RouteWs))
*/
func listenAndServeHTTPMux(addr string, mux http.Handler, rTimeout int, wTimeout int) error {
	srv := http.Server{
		Handler:      mux,
		Addr:         addr,
		ReadTimeout:  time.Duration(rTimeout) * time.Second,
		WriteTimeout: time.Duration(wTimeout) * time.Second,
	}
	return srv.ListenAndServe()
}
