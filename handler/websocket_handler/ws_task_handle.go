package websocket_handler

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type WsTaskHandler interface {
	ServeWs(c chan []byte, msg interface{}, conn interface{})
}

type WsTaskFunc func(c chan []byte, msg interface{}, conn interface{})

func (t WsTaskFunc) ServeWs(c chan []byte, msg interface{}, conn interface{}) {
	t(c, msg, conn)
}

type wsTaskHandle struct {
	handler WsTaskHandler
	timeOut time.Duration
}

var (
	wsHandlePoolMu sync.Mutex
	wsHandlePool   = make(map[string]*wsTaskHandle)
)

func RegisterWsTaskHandle(pattern string, handler WsTaskHandler, timeOut time.Duration) {
	wsHandlePoolMu.Lock()
	defer wsHandlePoolMu.Unlock()
	newHandle := &wsTaskHandle{
		handler: handler,
		timeOut: timeOut,
	}
	wsHandlePool[pattern] = newHandle
}

func GetWsTaskHandle(pattern string) (*wsTaskHandle, error) {
	wsHandlePoolMu.Lock()
	defer wsHandlePoolMu.Unlock()
	if handle, ok := wsHandlePool[pattern]; ok {
		return handle, nil
	} else {
		return nil, errors.New("can't find  handle")
	}
}

func DumpWsTaskHandle() {
	wsHandlePoolMu.Lock()
	defer wsHandlePoolMu.Unlock()
	for k, v := range wsHandlePool {
		fmt.Println(k, v)
	}
}
