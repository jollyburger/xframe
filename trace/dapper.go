package trace

import (
	"context"
	"encoding/json"
	"net"
	"sync/atomic"
	"time"
	"xframe/config"
	"xframe/log"
)

var (
	SAMPLE         uint32 = 0
	DEFAULT_SAMPLE uint32 = 100
)

type XContext interface {
	SetTaskId(int32) XContext
	SetSessionNo(string) XContext
	SetFuncName(string) XContext
	SetTraceFlag(bool) XContext
	SetDataMap(map[string]interface{}) XContext
	SetTime(string) XContext
	SetSpanId(int32) XContext
	GetTraceId() int32
	GetSpanId() int32
	GetSessionNo() string
	SendTraceData()
}

type Trace struct {
	context.Context                        //golib context
	ParentId        int32                  //ancestor trace id, send to traceaccess
	TraceId         int32                  //trace id, send to traceaccess
	SpanId          int32                  //span id: current trace id after rpc
	TaskId          int32                  //task id
	SessionNo       string                 //session no
	FuncName        string                 //function name
	TraceFlag       bool                   //trace switch
	Time            string                 //trace data start
	DataMap         map[string]interface{} //customized data
}

func InitTrace(session_no string, trace_id int32, span_id int32) *Trace {
	var sampling_base uint32
	var trace_flag bool
	sampling, err := config.GetConfigByKey("trace.sampling")
	if err != nil {
		sampling_base = DEFAULT_SAMPLE
	} else {
		sampling_base = uint32(sampling.(float64))
	}
	atomic.AddUint32(&SAMPLE, 1)
	if sampling_base == 0 {
		trace_flag = false
	} else if SAMPLE%sampling_base == 0 {
		trace_flag = true
	}
	root, err := config.GetConfigByKey("trace.root")
	if err == nil && root.(bool) {
		return newTrace(-1, -1, session_no, trace_flag)
	}
	return newTrace(trace_id, span_id, session_no, trace_flag)
}

func newTrace(trace_id, span_id int32, session_no string, trace_flag bool) *Trace {
	trace := new(Trace)
	trace.ParentId = trace_id
	trace.TraceId = span_id + 1
	trace.SpanId = span_id + 1
	trace.SessionNo = session_no
	trace.Time = time.Now().String()
	trace.TraceFlag = trace_flag
	return trace
}

func (t *Trace) SetTaskId(task_id int32) XContext {
	t.TaskId = task_id
	return t
}

func (t *Trace) SetSessionNo(session_no string) XContext {
	t.SessionNo = session_no
	return t
}

func (t *Trace) SetFuncName(func_name string) XContext {
	t.FuncName = func_name
	return t
}

func (t *Trace) SetTraceFlag(flag bool) XContext {
	t.TraceFlag = flag
	return t
}

func (t *Trace) SetDataMap(dm map[string]interface{}) XContext {
	t.DataMap = dm
	return t
}

func (t *Trace) SetTime(pt string) XContext {
	t.Time = pt
	return t
}

func (t *Trace) SetSpanId(span_id int32) XContext {
	t.SpanId = span_id
	return t
}

func (t *Trace) GetTraceId() int32 {
	return t.TraceId
}

func (t *Trace) GetSpanId() int32 {
	return t.SpanId
}

func (t *Trace) GetSessionNo() string {
	return t.SessionNo
}

//UDP client
//TODO: udp client pool, MOVE TO utils
func sendUDPData(b []byte, addr string) error {
	udp_addr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	udp_conn, err := net.DialUDP("udp", nil, udp_addr)
	if err != nil {
		return err
	}
	_, err = udp_conn.Write(b)
	return err
}

func (t *Trace) SendTraceData() {
	if !t.TraceFlag {
		return
	}
	//get trace addr
	trace_addr, err := config.GetConfigByKey("trace.addr")
	if err != nil {
		log.ERRORF("send trace data error: %v", err)
		return
	}
	//encoding data in json type
	b, err := json.Marshal(t)
	if err != nil {
		log.ERRORF("marshal json data error: %v", err)
		return
	}
	//send udp data
	err = sendUDPData(b, trace_addr.(string))
	if err != nil {
		log.ERRORF("send trace data error: %v", err)
		return
	}
}
