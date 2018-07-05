package logic

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"xframe/handler/http_handler"
	"xframe/log"
	"xframe/server"
	"xframe/trace"
	"xframe/utils"
)

func init() {
	//register http handler
	http_handler.RegisterHTTPTaskHandle("Echo", http_handler.XHTTPTaskFunc(echo_serve), 20*time.Second)
	server.RouteHTTP = RouteHTTP
}

func RouteHTTP(w http.ResponseWriter, r *http.Request) {
	log.DEBUG(r)
	r.ParseForm()
	//get session_no in header
	session_no := r.Header.Get("X-Session-No")
	if session_no == "" {
		session_no = utils.NewUUIDV4().String()
	}
	//get trace_id and span_id in header
	var trace_id, span_id int
	ids := r.Header.Get("X-Trace-Context")
	if ids != "" && len(strings.Split(ids, ":")) == 2 {
		trace_id, _ = strconv.Atoi(strings.Split(ids, ":")[0])
		span_id, _ = strconv.Atoi(strings.Split(ids, ":")[1])
	} else {
		trace_id, span_id = -1, -1
	}
	//add customized data in ctx
	ctx := trace.InitTrace(session_no, int32(trace_id), int32(span_id))
	go ctx.SendTraceData()
	//process input value
	action := r.FormValue("Action")
	if action == "" {
		log.ERROR("can not find action")
		buf, _ := FormatResponse("On Data In", -99, "can not find action", "")
		server.SendHTTPResponse(ctx, w, buf)
		return
	}
	new_task, err := http_handler.NewHTTPTask(action)
	if err != nil {
		log.DEBUG("new task fail, ", err)
		return
	}
	_, err = new_task.Run(ctx, w, r)
	if err != nil {
		log.ERRORF("run task(%s) fail:%s", new_task.Id, err)
		return
	}
}
