package handler_common

import (
	"reflect"
	"runtime"
	"strings"
	"time"
)

type TaskState int

const (
	StateNew TaskState = iota
	StateRun
	StateFinished
)

var (
	GlobalTaskId int32
)

type CommonTask struct {
	Id       int32
	Gid      uint64
	FuncName string
	Time     time.Duration
	State    TaskState
}

// through handler, get the handler name
func GetTaskFuncName(taskHandler interface{}) string {
	funcInfo := runtime.FuncForPC(reflect.ValueOf(taskHandler).Pointer()).Name()
	return strings.Split(funcInfo, ".")[1]
}
