package http_handler

import (
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"xframe/handler/handler_common"
	"xframe/trace"
	"xframe/utils"
)

var (
	httpTaskPoolMu sync.Mutex
	httpTaskPool   = make(map[int32]*HTTPTask)
)

type HTTPTask struct {
	handler_common.CommonTask
	Pattern    string
	Handler    XHttpTaskHandler
	isFinished chan bool
}

func NewHTTPTask(pattern string) (task *HTTPTask, err error) {
	taskHandle, err := GetHTTPTaskHandle(pattern)
	if err != nil {
		return
	}
	task = &HTTPTask{
		CommonTask: handler_common.CommonTask{
			Id:    atomic.AddInt32(&handler_common.GlobalTaskId, 1),
			State: handler_common.StateNew,
			Time:  taskHandle.timeOut,
		},
		Pattern:    pattern,
		Handler:    taskHandle.handler,
		isFinished: make(chan bool),
	}
	httpTaskPoolMu.Lock()
	httpTaskPool[task.Id] = task
	httpTaskPoolMu.Unlock()
	return
}

func (t *HTTPTask) Run(ctx trace.XContext, rw http.ResponseWriter, r *http.Request) (res []byte, err error) {
	t.setState(handler_common.StateRun)
	funcName := handler_common.GetTaskFuncName(t.Handler)
	go func() {
		t.Gid = utils.GetGID()
		t.FuncName = funcName
		if utils.CheckWrapPanic() {
			utils.GoSafeHTTP(ctx, rw, r, t.Handler.ServeXHTTP)
		} else {
			t.Handler.ServeXHTTP(ctx, utils.CtxInHttpRspHeader(ctx, rw), utils.CtxInHttpReqHeader(ctx, r))
		}
		t.isFinished <- true
	}()
	if t.Time > 0 {
		select {
		case <-t.isFinished:
			t.setState(handler_common.StateFinished)
		case <-time.After(t.Time):
			t.setState(handler_common.StateFinished)
			err = errors.New("task time out")
		}
	} else {
		select {
		case <-t.isFinished:
			t.setState(handler_common.StateFinished)
		}
	}
	httpTaskPoolMu.Lock()
	delete(httpTaskPool, t.Id)
	httpTaskPoolMu.Unlock()
	return
}

func (t *HTTPTask) setState(state handler_common.TaskState) {
	t.State = state
}

func LenHTTPTasks() int {
	return len(httpTaskPool)
}

func GetHTTPTaskByGid(gid uint64) (task interface{}) {
	httpTaskPoolMu.Lock()
	defer httpTaskPoolMu.Unlock()
	for _, t := range httpTaskPool {
		if t.Gid == gid {
			return t
		}
	}
	return nil
}

func DumpHTTPTasks() (tasks map[int32]*HTTPTask) {
	tasks = make(map[int32]*HTTPTask)
	httpTaskPoolMu.Lock()
	defer httpTaskPoolMu.Unlock()
	for k, v := range httpTaskPool {
		tasks[k] = v
	}
	return
}
