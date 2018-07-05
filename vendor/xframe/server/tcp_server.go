package server

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"xframe/log"
)

type tcpServer struct {
	listener net.Listener
	// About connection
	maxConnnectionId uint64
	connections      map[uint64]*TcpConnection
	connectionMutex  sync.Mutex
	// About server start and stop
	stopFlag int32
	stopWait sync.WaitGroup
}

func newTcpServer(listener net.Listener) *tcpServer {
	return &tcpServer{
		listener:    listener,
		connections: make(map[uint64]*TcpConnection),
	}
}

func (s *tcpServer) serve() (err error) {
	log.DEBUGF("server start service at %s", s.listener.Addr())
	for {
		c, err := s.listener.Accept()
		if err != nil {
			log.ERRORF("Accept fail:%v", err)
			fmt.Println("Accept fail:", err)
			break
		}
		connection, err := s.newConnection(c)
		log.DEBUGF("new server connection [ %s -> %s ]", c.RemoteAddr(), c.LocalAddr())
		OnConnect(connection)
		go s.serveConnection(connection)
	}
	return
}

func (s *tcpServer) newConnection(conn net.Conn) (c *TcpConnection, err error) {
	c = newTcpConnection(conn)
	if err = c.SetKeepAlive(defaultKeepAlivePeriod); err != nil {
		return
	}
	s.addConnection(c)
	return
}

func (s *tcpServer) addConnection(c *TcpConnection) {
	s.connectionMutex.Lock()
	defer s.connectionMutex.Unlock()
	s.connections[c.id] = c
	s.stopWait.Add(1)
}

func (s *tcpServer) serveConnection(c *TcpConnection) {
	for {
		req, err := c.Receive()
		if err != nil {
			log.DEBUGF("connection [ %s -> %s ] is closed", c.conn.RemoteAddr(), c.conn.LocalAddr())
			s.delConnection(c)
			return
		}
		go OnDataIn(c, req)
	}
}

func (s *tcpServer) lenConnection() int {
	return len(s.connections)
}

func (s *tcpServer) stop() bool {
	if atomic.CompareAndSwapInt32(&s.stopFlag, 0, 1) {
		s.listener.Close()
		s.closeConnections()
		s.stopWait.Wait()
		return true
	}
	return false
}

func (s *tcpServer) delConnection(c *TcpConnection) {
	s.connectionMutex.Lock()
	defer s.connectionMutex.Unlock()
	delete(s.connections, c.id)
	s.stopWait.Done()
}

func (s *tcpServer) closeConnections() {
	for _, connection := range s.connections {
		connection.Close()
	}
}
