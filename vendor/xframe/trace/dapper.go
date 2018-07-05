package trace

import (
	"context"
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var (
	GTracer opentracing.Tracer
)

type XContext interface {
	context.Context
	/*SetTaskId(int32) XContext
	SetSessionNo(string) XContext
	SetFuncName(string) XContext
	SetTraceFlag(bool) XContext
	SetDataMap(map[string]interface{}) XContext
	SetTime(string) XContext
	SetSpanId(int32) XContext
	GetTraceId() int32
	GetSpanId() int32
	GetSessionNo() string
	SendTraceData()*/
}

type Trace struct {
	/*ParentId  int32                  `json:"parent_id"`  //ancestor trace id, send to traceaccess
	TraceId   int32                  `json:"trace_id"`   //trace id, send to traceaccess
	SpanId    int32                  `json:"span_id"`    //span id: current trace id after rpc
	TaskId    int32                  `json:"task_id"`    //task id
	SessionNo string                 `json:"session_no"` //session no
	FuncName  string                 `json:"func_name"`  //function name
	TraceFlag bool                   `json:"-"`          //trace switch
	Time      string                 `json:"time"`       //trace data start
	DataMap   map[string]interface{} `json:"data_map"`   //customized data*/
}

func InitContext() context.Context {
	/*var sampling_base uint32
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
	return newTrace(trace_id, span_id, session_no, trace_flag)*/
	return context.Background()
}

/*
 * service_name: tracing name in UI
 * trace_type: const, probabilistic, rateLimiting, or remote(deprecated)
 * params:
 * - for const, 0/1 means switch on/off
 * - for probabilistic, from 0 to 1
 * - for rateLimiting, the number of spans per second
 */

//TODO: more flexible configuration
func InitTracer(service_name string, trace_type string, params float64) (err error) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  trace_type,
			Param: params,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 3 * time.Second,
		},
	}
	GTracer, _, err = cfg.New(service_name)
	return
}

func ExtractHTTPTracer(header http.Header) (opentracing.SpanContext, error) {
	return GTracer.Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(header))
}

/*
 * trace: TraceId in Protocol Buffer's Head
 */
func ExtractBinaryTracer(trace string) (opentracing.SpanContext, error) {
	return jaeger.ContextFromString(trace)
}

func StartSpan(operation_name string, ext_ctx opentracing.SpanContext) opentracing.Span {
	if ext_ctx != nil {
		return GTracer.StartSpan(operation_name, opentracing.ChildOf(ext_ctx))
	}
	return GTracer.StartSpan(operation_name)
}

func FinishSpan(sp opentracing.Span) {
	sp.Finish()
}

func ContextWithSpan(ctx context.Context, sp opentracing.Span) context.Context {
	return opentracing.ContextWithSpan(context.Context(ctx), sp)
}

func SpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}

func SetSpanTag(sp opentracing.Span, key string, value interface{}) {
	sp.SetTag(key, value)
}

func SpanToSpanContext(span opentracing.Span) jaeger.SpanContext {
	return span.Context().(jaeger.SpanContext)
}

/*func newTrace(trace_id, span_id int32, session_no string, trace_flag bool) *Trace {
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
}*/

//UDP client
//TODO: udp client pool, MOVE TO utils
/*func sendUDPData(b []byte, addr string) error {
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
}*/
