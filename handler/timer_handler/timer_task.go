package timer_handler

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/jollyburger/xframe/handler/handler_common"
	"github.com/jollyburger/xframe/log"
	"github.com/jollyburger/xframe/utils"
)

var (
	timerTaskPoolMu    sync.Mutex
	timerTaskPool                 = make(map[int32]*TimerTask)
	timerTaskServeOnce *sync.Once = &sync.Once{}
)

type TimerTask struct {
	handler_common.CommonTask
	Type       int32
	Handler    TimerTaskHandler
	singleton  bool
	timer      *time.Ticker
	isFinished chan bool
}

func newTimerTask(tType int32, handle *timerTaskHandle) (task *TimerTask, err error) {
	task = &TimerTask{
		CommonTask: handler_common.CommonTask{
			Id:    atomic.AddInt32(&handler_common.GlobalTaskId, 1),
			State: handler_common.StateNew,
			Time:  handle.intervalTime,
		},
		Type:       tType,
		Handler:    handle.handler,
		singleton:  handle.singleton,
		isFinished: make(chan bool),
	}
	timerTaskPoolMu.Lock()
	timerTaskPool[task.Id] = task
	timerTaskPoolMu.Unlock()
	return
}

func (t *TimerTask) Run() {
	funcName := handler_common.GetTaskFuncName(t.Handler)
	go func() {
		t.Gid = utils.GetGID()
		t.FuncName = funcName
		if utils.CheckWrapPanic() {
			utils.GoSafeTimer(t.Handler.ServeTimer)
		} else {
			t.Handler.ServeTimer()
		}
		t.isFinished <- true
	}()
	<-t.isFinished
	timerTaskPoolMu.Lock()
	delete(timerTaskPool, t.Id)
	timerTaskPoolMu.Unlock()
	return
}

func (t *TimerTask) setState(state handler_common.TaskState) {
	t.State = state
}

func TimerTaskServe() {
	timerTaskServeOnce.Do(timerTaskServe)
}

func timerTaskServe() {
	timerHandlePoolMu.Lock()
	defer timerHandlePoolMu.Unlock()
	for tType, handle := range timerHandlePool {
		//run first
		task, err := newTimerTask(tType, handle)
		if err != nil {
			log.ERRORF("create timer task[%d] fail", tType)
			continue
		}
		go task.Run()
		//set timer
		taskTimer := time.NewTicker(handle.intervalTime)
		go runTimerTask(tType, handle, taskTimer)
	}
}

func runTimerTask(tType int32, handle *timerTaskHandle, t *time.Ticker) {
	for _ = range t.C {
		// check singleton
		isRunning := false
		timerTaskPoolMu.Lock()
		for _, rtask := range timerTaskPool {
			if rtask.Type == tType {
				if rtask.singleton {
					isRunning = true
					break
				}
			}
		}
		timerTaskPoolMu.Unlock()
		if isRunning {
			continue
		}
		task, err := newTimerTask(tType, handle)
		if err != nil {
			log.ERRORF("create timer task fail, type(%d)", tType)
			continue
		}
		go task.Run()
	}
}

func LenTimerTasks() int {
	return len(timerTaskPool)
}

func GetTimerTaskByGid(gid uint64) (task interface{}) {
	timerTaskPoolMu.Lock()
	defer timerTaskPoolMu.Unlock()
	for _, t := range timerTaskPool {
		if t.Gid == gid {
			return t
		}
	}
	return nil
}

func DumpTimerTasks() (tasks map[int32]*TimerTask) {
	tasks = make(map[int32]*TimerTask)
	timerTaskPoolMu.Lock()
	defer timerTaskPoolMu.Unlock()
	for k, v := range timerTaskPool {
		tasks[k] = v
	}
	return
}
