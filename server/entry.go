package server

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"xframe/trace"
)

// tcp
func RunTCP(listen_addr string, listen_port int) (err error) {
	var listener net.Listener
	listen_ip, err := parseListenAddr(listen_addr)
	if err != nil {
		return
	}
	address := net.JoinHostPort(listen_ip, strconv.Itoa(listen_port))
	if listener, err = net.Listen("tcp", address); err != nil {
		return
	}
	server := newTcpServer(listener)
	server.serve()
	return
}

func SetTCPClientConnLimit(limit int) {
	setClientConnectionLimit(limit)
}

// Request with response
func SendTCPRequest(s_peer_addr string, i_peer_port int, req []byte, timeOut uint32) (res []byte, err error) {
	return sendClientRequest(s_peer_addr, i_peer_port, req, timeOut)
}

// Request without response
func SendTCPRequestNoResponse(s_peer_addr string, i_peer_port int, req []byte, timeOut uint32) (err error) {
	return sendClientRequestNoResponse(s_peer_addr, i_peer_port, req, timeOut)
}

func SendTCPResponse(connection *TcpConnection, res []byte) (err error) {
	_, err = connection.Send(res)
	return
}

func ParseListenAddr(listen_addr string) (listen_ip string, err error) {
	return parseListenAddr(listen_addr)
}

func IsIPv4(ip string) bool {
	return isIPv4(ip)
}

func NewTcpConnection(conn net.Conn) *TcpConnection {
	return newTcpConnection(conn)
}

// ========================================================================
// http
func RunHTTP(listen_addr string, listen_port int) (err error) {
	listen_ip, err := parseListenAddr(listen_addr)
	if err != nil {
		return
	}
	address := net.JoinHostPort(listen_ip, strconv.Itoa(listen_port))
	err = listenAndServeHTTP(address)
	return
}

//customized http multiplexer
func RunHTTPMux(listen_addr string, listen_port int, mux http.Handler) (err error) {
	listen_ip, err := parseListenAddr(listen_addr)
	if err != nil {
		return
	}
	address := net.JoinHostPort(listen_ip, strconv.Itoa(listen_port))
	err = listenAndServeHTTPMux(address, mux)
	return
}

func SendHTTPRequest(ctx trace.XContext, uri string, params map[string]interface{}, timeOut uint32) (res []byte, err error) {
	return sendHttpRequest(ctx, uri, params, timeOut)
}

func SendHTTPPostRequest(ctx trace.XContext, uri string, body_type string, body io.Reader, timeOut uint32) (res []byte, err error) {
	return sendHttpPostRequest(ctx, uri, body_type, body, timeOut)
}

func SendHTTPMethodRequest(ctx trace.XContext, method, uri string, body io.Reader, timeOut uint32) (res []byte, err error) {
	return sendHttpMethodRequest(ctx, method, uri, body, timeOut)
}

func SendHTTPRequestByLimit(ctx trace.XContext, method, urlStr string, reader io.Reader, header http.Header, rate float64, capacity int64) ([]byte, error) {
	return sendHttPRequestBylimit(ctx, method, urlStr, reader, header, rate, capacity)
}

func SendHTTPResponse(ctx trace.XContext, w http.ResponseWriter, b []byte) error {
	_, err := w.Write(b)
	return err
}

//ratelimit
func InitRateLimit(rate float64, capacity int64) error {
	if rate <= 0.0 || capacity <= 0 {
		return errors.New("ratelimiter input error")
	}
	defaultCapacity = capacity
	defaultRate = rate
	return nil
}
