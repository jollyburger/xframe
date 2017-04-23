package websocket_handler

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
	"xframe/handler/handler_common"
	"xframe/utils"
)

var (
	wsTaskPoolMu sync.Mutex
	wsTaskPool   = make(map[int32]*WsTask)
)

type WsTask struct {
	handler_common.CommonTask
	Pattern string
	Handler WsTaskHandler
	msgChan chan []byte
}

func NewWsTask(pattern string) (task *WsTask, err error) {
	taskHandle, err := GetWsTaskHandle(pattern)
	if err != nil {
		return
	}
	task = &WsTask{
		CommonTask: handler_common.CommonTask{
			Id:    atomic.AddInt32(&handler_common.GlobalTaskId, 1),
			State: handler_common.StateNew,
			Time:  taskHandle.timeOut,
		},
		Pattern: pattern,
		Handler: taskHandle.handler,
		msgChan: make(chan []byte),
	}
	wsTaskPoolMu.Lock()
	wsTaskPool[task.Id] = task
	wsTaskPoolMu.Unlock()
	return
}

func (t *WsTask) Run(req interface{}, conn interface{}) (res []byte, err error) {
	t.setState(handler_common.StateRun)
	funcName := handler_common.GetTaskFuncName(t.Handler)
	var ok bool
	go func() {
		t.Gid = utils.GetGID()
		t.FuncName = funcName
		t.Handler.ServeWs(t.msgChan, req, conn)
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
	wsTaskPoolMu.Lock()
	delete(wsTaskPool, t.Id)
	wsTaskPoolMu.Unlock()
	return
}

func (t *WsTask) setState(state handler_common.TaskState) {
	t.State = state
}

func LenWsTasks() int {
	return len(wsTaskPool)
}

func GetWsTaskByGid(gid uint64) (task interface{}) {
	wsTaskPoolMu.Lock()
	defer wsTaskPoolMu.Unlock()
	for _, t := range wsTaskPool {
		if t.Gid == gid {
			return t
		}
	}
	return nil
}

func DumpWsTasks() (tasks map[int32]*WsTask) {
	tasks = make(map[int32]*WsTask)
	wsTaskPoolMu.Lock()
	defer wsTaskPoolMu.Unlock()
	for k, v := range wsTaskPool {
		tasks[k] = v
	}
	return
}
