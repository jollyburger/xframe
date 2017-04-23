package tcp_handler

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
	"xframe/handler/handler_common"
	"xframe/trace"
	"xframe/utils"
)

var (
	tcpTaskPoolMu sync.Mutex
	tcpTaskPool   = make(map[int32]*TCPTask)
)

type TCPTask struct {
	handler_common.CommonTask
	Type    int32
	Handler TCPTaskHandler
	msgChan chan []byte
}

func NewTCPTask(tType int32) (task *TCPTask, err error) {
	taskHandle, err := GetTCPTaskHandle(tType)
	if err != nil {
		return
	}
	task = &TCPTask{
		CommonTask: handler_common.CommonTask{
			Id:    atomic.AddInt32(&handler_common.GlobalTaskId, 1),
			State: handler_common.StateNew,
			Time:  taskHandle.timeOut,
		},
		Type:    tType,
		Handler: taskHandle.handler,
		msgChan: make(chan []byte),
	}
	tcpTaskPoolMu.Lock()
	tcpTaskPool[task.Id] = task
	tcpTaskPoolMu.Unlock()
	return
}

func (t *TCPTask) Run(ctx trace.XContext, req interface{}) (res []byte, err error) {
	t.setState(handler_common.StateRun)
	funcName := handler_common.GetTaskFuncName(t.Handler)
	var ok bool
	go func() {
		t.Gid = utils.GetGID()
		t.FuncName = funcName
		if utils.CheckWrapPanic() {
			utils.GoSafeTCP(ctx, t.msgChan, req, t.Handler.ServeTCP)
		} else {
			t.Handler.ServeTCP(ctx, t.msgChan, req)
		}
	}()
	if t.Time > 0 {
		select {
		case res, ok = <-t.msgChan:
			if !ok {
				err = errors.New("task fail, close")
			}
			t.setState(handler_common.StateFinished)
		case <-time.After(t.Time):
			t.setState(handler_common.StateFinished)
			err = errors.New("task time out")

		}
	} else {
		select {
		case res, ok = <-t.msgChan:
			if !ok {
				err = errors.New("task fail, close")
			}
			t.setState(handler_common.StateFinished)
		}
	}
	tcpTaskPoolMu.Lock()
	delete(tcpTaskPool, t.Id)
	tcpTaskPoolMu.Unlock()
	return
}

func (t *TCPTask) setState(state handler_common.TaskState) {
	t.State = state
}

func LenTCPTasks() int {
	return len(tcpTaskPool)
}

func GetTCPTaskByGid(gid uint64) (task interface{}) {
	tcpTaskPoolMu.Lock()
	defer tcpTaskPoolMu.Unlock()
	for _, t := range tcpTaskPool {
		if t.Gid == gid {
			return t
		}
	}
	return nil
}

func DumpTCPTasks() (tasks map[int32]*TCPTask) {
	tasks = make(map[int32]*TCPTask)
	tcpTaskPoolMu.Lock()
	defer tcpTaskPoolMu.Unlock()
	for k, v := range tcpTaskPool {
		tasks[k] = v
	}
	return
}
