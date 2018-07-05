package timer_handler

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type TimerTaskHandler interface {
	ServeTimer()
}

type TimerTaskFunc func()

func (t TimerTaskFunc) ServeTimer() {
	t()
}

type timerTaskHandle struct {
	handler      TimerTaskHandler
	intervalTime time.Duration
	singleton    bool
}

var (
	timerHandlePoolMu sync.Mutex
	timerHandlePool   = make(map[int32]*timerTaskHandle)
)

func RegisterTimerTaskHandle(id int32, handler TimerTaskHandler, intervalTime time.Duration, singleton bool) {
	timerHandlePoolMu.Lock()
	defer timerHandlePoolMu.Unlock()
	newHandle := &timerTaskHandle{
		handler:      handler,
		intervalTime: intervalTime,
		singleton:    singleton,
	}
	timerHandlePool[id] = newHandle
}

func GetTimerTaskHandle(id int32) (*timerTaskHandle, error) {
	timerHandlePoolMu.Lock()
	defer timerHandlePoolMu.Unlock()
	if handle, ok := timerHandlePool[id]; ok {
		return handle, nil
	} else {
		return nil, errors.New("can't not find handle")
	}
}

func DumpTimerTaskHandle() {
	timerHandlePoolMu.Lock()
	defer timerHandlePoolMu.Unlock()
	for k, v := range timerHandlePool {
		fmt.Println(k, v)
	}
}
