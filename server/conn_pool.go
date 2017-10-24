package server

import (
	"errors"
	"net"
	"sync"
	"time"
	"xframe/log"
)

var (
	connectionLimit       int = 100
	defaultConnectTimeout     = 3 * time.Second
)

func setClientConnectionLimit(limit int) {
	connectionLimit = limit
}

//as client
type clientTcpConnection struct {
	isOldConn bool
	tcpConn   *TcpConnection
}

func (ctc *clientTcpConnection) IsOld() bool {
	return ctc.isOldConn
}

func (ctc *clientTcpConnection) LocalAddr() string {
	return ctc.tcpConn.Conn().LocalAddr().String()
}

func (ctc *clientTcpConnection) RemoteAddr() string {
	return ctc.tcpConn.Conn().RemoteAddr().String()
}

func (ctc *clientTcpConnection) Close() {
	ctc.tcpConn.Close()
}

func (ctc *clientTcpConnection) Send(req []byte) (int, error) {
	return ctc.tcpConn.Send(req)
}

func (ctc *clientTcpConnection) Receive() ([]byte, error) {
	return ctc.tcpConn.Receive()
}

func (ctc *clientTcpConnection) SetDeadline(t time.Duration) error {
	return ctc.tcpConn.SetDeadline(t)
}

//==================================
type ConnPool struct {
	mu      sync.RWMutex
	conns   chan *clientTcpConnection
	factory Factory
}

type Factory func() (*clientTcpConnection, error)

func Adaptor(network, addr string, timeout time.Duration) Factory {
	ff := func() (*clientTcpConnection, error) {
		if timeout > defaultConnectTimeout {
			timeout = defaultConnectTimeout
		}
		conn, err := net.DialTimeout(network, addr, timeout)
		if err != nil {
			return nil, err
		}
		log.DEBUGF("new client connection [%s -> %s]", conn.LocalAddr(), conn.RemoteAddr())
		tcpConn := newTcpConnection(conn)
		err = tcpConn.SetKeepAlive(defaultKeepAlivePeriod)
		if err != nil {
			return nil, err
		}
		client_tcp_conn := &clientTcpConnection{
			tcpConn: tcpConn,
		}
		return client_tcp_conn, nil
	}
	return ff
}

func newConnPool(network, addr string, timeout time.Duration) *ConnPool {
	conn_pool := new(ConnPool)
	conn_pool.conns = make(chan *clientTcpConnection, connectionLimit)
	conn_pool.factory = Adaptor(network, addr, timeout)
	return conn_pool
}

func (cp *ConnPool) Get() (*clientTcpConnection, error) {
	conns := cp.getConns()
	if conns == nil {
		return nil, errors.New("connection pool is empty")
	}
	select {
	case conn := <-conns:
		if conn != nil {
			return conn, nil
		}
		return nil, errors.New("get connection from pool error")
	default:
		conn, err := cp.factory()
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}

func (cp *ConnPool) getConns() chan *clientTcpConnection {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	return cp.conns
}

func (cp *ConnPool) Put(conn *clientTcpConnection) {
	conns := cp.getConns()
	if conns == nil {
		conns = make(chan *clientTcpConnection, connectionLimit)
		cp.conns = conns
	}
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if len(conns) >= connectionLimit {
		conn.tcpConn.Close()
	} else {
		conn.isOldConn = true
		conns <- conn
	}
}

func (cp *ConnPool) Len() int {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	return len(cp.conns)
}

func (cp *ConnPool) CloseAll() {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if len(cp.conns) != 0 {
		log.DEBUG("closing connection")
		for c := range cp.conns {
			c.Close()
		}
	}
	close(cp.conns)
}
