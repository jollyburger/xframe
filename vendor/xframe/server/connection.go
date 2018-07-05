package server

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type TcpConnection struct {
	id   uint64
	conn net.Conn
	// About send and receive
	reader    *reader
	writer    *writer
	sendMutex sync.Mutex
	recvMutex sync.Mutex
	// About close
	closeFlag int32
}

var (
	globalTcpConnectionId uint64
	noDeadline            = time.Time{}
)

const (
	defaultKeepAlivePeriod = 10 * time.Second
)

func newTcpConnection(conn net.Conn) *TcpConnection {
	return &TcpConnection{
		id:     atomic.AddUint64(&globalTcpConnectionId, 1),
		conn:   conn,
		reader: newReader(conn),
		writer: newWriter(conn),
	}
}

func (c *TcpConnection) Id() uint64     { return c.id }
func (c *TcpConnection) Conn() net.Conn { return c.conn }
func (c *TcpConnection) IsClosed() bool { return atomic.LoadInt32(&c.closeFlag) != 0 }

func (c *TcpConnection) Receive() (msg []byte, err error) {
	c.recvMutex.Lock()
	defer c.recvMutex.Unlock()
	if msg, err = c.reader.readPacket(); err != nil {
		c.Close()
	}
	return
}

func (c *TcpConnection) Send(msg []byte) (n int, err error) {
	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()
	if n, err = c.writer.writePacket(msg); err != nil {
		c.Close()
	}
	//DataOut hook
	OnDataOut(c.conn, msg)
	return
}

func (c *TcpConnection) Close() {
	if atomic.CompareAndSwapInt32(&c.closeFlag, 0, 1) {
		c.conn.Close()
		//Connection hook
		OnDisconnect(c)
	}
}

func (c *TcpConnection) SetKeepAlive(period time.Duration) (err error) {
	if tc, ok := c.conn.(*net.TCPConn); ok {
		if err = tc.SetKeepAlive(true); err != nil {
			return
		}
		if err = tc.SetKeepAlivePeriod(period); err != nil {
			return
		}
	}
	return
}

func (c *TcpConnection) SetReadDeadline(timeOut time.Duration) (err error) {
	if tc, ok := c.conn.(*net.TCPConn); ok {
		if timeOut != 0 {
			err = tc.SetReadDeadline(time.Now().Add(timeOut))
		} else {
			err = tc.SetReadDeadline(noDeadline)
		}
	}
	return
}

func (c *TcpConnection) SetWriteDeadline(timeOut time.Duration) (err error) {
	if tc, ok := c.conn.(*net.TCPConn); ok {
		if timeOut != 0 {
			err = tc.SetWriteDeadline(time.Now().Add(timeOut))
		} else {
			err = tc.SetWriteDeadline(noDeadline)
		}
	}
	return
}

func (c *TcpConnection) SetDeadline(timeOut time.Duration) (err error) {
	if tc, ok := c.conn.(*net.TCPConn); ok {
		if timeOut != 0 {
			err = tc.SetDeadline(time.Now().Add(timeOut))
		} else {
			err = tc.SetDeadline(noDeadline)
		}
	}
	return
}
