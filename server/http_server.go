package server

import (
	"net/http"
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
func listenAndServeHTTPMux(addr string, mux http.Handler) error {
	return http.ListenAndServe(addr, mux)
}
