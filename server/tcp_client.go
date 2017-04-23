package server

import (
	"strconv"
	"sync"
	"time"
	"xframe/log"
)

const (
	maxBadConnRetries = 2
)

var (
	clientTcpConnectionPoolMu sync.Mutex
	clientTcpConnectionPool   = make(map[string]*ConnPool)
)

func newClientTcpConnection(s_peer_addr string, i_peer_port int, timeOut uint32) (c *clientTcpConnection, err error) {
	remote_addr := s_peer_addr + ":" + strconv.Itoa(i_peer_port)
	clientTcpConnectionPoolMu.Lock()
	connPool, ok := clientTcpConnectionPool[remote_addr]
	if ok {
		c, err = connPool.Get()
		log.DEBUGF("get client connection [ %s -> %s ] from pool", c.LocalAddr(), c.RemoteAddr())
		clientTcpConnectionPoolMu.Unlock()
		return
	}
	clientTcpConnectionPoolMu.Unlock()
	return connectServer("tcp", remote_addr, time.Duration(timeOut)*time.Second)
}

func connectServer(network, address string, timeout time.Duration) (c *clientTcpConnection, err error) {
	conn_pool := newConnPool(network, address, timeout)
	clientTcpConnectionPoolMu.Lock()
	defer clientTcpConnectionPoolMu.Unlock()
	clientTcpConnectionPool[address] = conn_pool
	c, err = conn_pool.Get()
	return
}

func freeClientTcpConnection(s_peer_addr string, i_peer_port int, c *clientTcpConnection) {
	remote_addr := s_peer_addr + ":" + strconv.Itoa(i_peer_port)
	clientTcpConnectionPoolMu.Lock()
	defer clientTcpConnectionPoolMu.Unlock()
	conn_pool, ok := clientTcpConnectionPool[remote_addr]
	if !ok {
		c.Close()
		return
	}
	connNum := conn_pool.Len()
	if connNum >= connectionLimit {
		c.Close()
	} else {
		conn_pool.Put(c)
		clientTcpConnectionPool[remote_addr] = conn_pool
	}
}

func closeClientTcpConnection(s_peer_addr string, i_peer_port int) {
	remote_addr := s_peer_addr + ":" + strconv.Itoa(i_peer_port)
	clientTcpConnectionPoolMu.Lock()
	defer clientTcpConnectionPoolMu.Unlock()
	conn_pool, ok := clientTcpConnectionPool[remote_addr]
	if ok {
		conn_pool.CloseAll()
		delete(clientTcpConnectionPool, remote_addr)
	}
}

func sendClientRequest(s_peer_addr string, i_peer_port int, req []byte, timeOut uint32) ([]byte, error) {
	var connection *clientTcpConnection
	var res []byte
	var err error
	// Retry
	for i := 0; i < maxBadConnRetries; i++ {
		connection, err = newClientTcpConnection(s_peer_addr, i_peer_port, timeOut)
		if err != nil {
			break
		}
		if err = connection.SetDeadline(time.Duration(timeOut) * time.Second); err != nil {
			if connection.IsOld() {
				closeClientTcpConnection(s_peer_addr, i_peer_port)
				continue
			}
			break
		}
		//send msg
		writeErrChan := make(chan error)
		go func() {
			_, e := connection.Send(req)
			writeErrChan <- e
		}()
		select {
		case err = <-writeErrChan:
			if err != nil {
				if connection.IsOld() {
					closeClientTcpConnection(s_peer_addr, i_peer_port)
					continue
				}
				break
			}
		case <-time.After(50 * time.Millisecond):
			err = <-writeErrChan
		}
		if err == nil {
			readErrChan := make(chan error)
			go func() {
				tmp_res, err := connection.Receive()
				res = make([]byte, 0, len(tmp_res))
				res = append(res, tmp_res...)
				//log.DEBUG("checking", len(res), res)
				readErrChan <- err
			}()
			select {
			case err = <-readErrChan:
				if err != nil {
					if connection.IsOld() {
						closeClientTcpConnection(s_peer_addr, i_peer_port)
						continue
					}
					break
				}
			case <-time.After(50 * time.Millisecond):
				err = <-readErrChan
			}
			break
		}
		break
	}
	//put back in conn pool
	if err == nil && connection != nil {
		freeClientTcpConnection(s_peer_addr, i_peer_port, connection)
		return res, err
	}
	// close connection
	if err != nil && connection != nil {
		connection.Close()
	}
	return res, err
}

func sendClientRequestNoResponse(s_peer_addr string, i_peer_port int, req []byte, timeOut uint32) (err error) {
	var connection *clientTcpConnection
	for i := 0; i < maxBadConnRetries; i++ {
		connection, err = newClientTcpConnection(s_peer_addr, i_peer_port, timeOut)
		if err != nil {
			return
		}
		if err = connection.SetDeadline(time.Duration(timeOut) * time.Second); err != nil {
			if connection.IsOld() {
				closeClientTcpConnection(s_peer_addr, i_peer_port)
				continue
			}
			break
		}
		writeErrChan := make(chan error)
		go func() {
			_, e := connection.Send(req)
			writeErrChan <- e
		}()
		select {
		case err = <-writeErrChan:
			if err != nil {
				if connection.IsOld() {
					closeClientTcpConnection(s_peer_addr, i_peer_port)
					continue
				}
				break
			}
		case <-time.After(50 * time.Millisecond):
			err = <-writeErrChan
		}
		break
	}
	if err == nil && connection != nil {
		freeClientTcpConnection(s_peer_addr, i_peer_port, connection)
		return
	}
	if err != nil && connection != nil {
		connection.Close()
	}
	return
}

func LenClientTcpConnections(s_peer_addr string, i_peer_port int) (plen int) {
	remote_addr := s_peer_addr + ":" + strconv.Itoa(i_peer_port)
	conn_pool, ok := clientTcpConnectionPool[remote_addr]
	if ok {
		plen = conn_pool.Len()
	}
	return
}
