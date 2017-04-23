package logic

import (
	"net/http"
	"xframe/log"
	"xframe/server"
	"xframe/trace"
)

func echo_serve(ctx trace.XContext, rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.ERROR("parse input error")
		buf, _ := FormatResponse("Echo", -1, "parse input error", "")
		server.SendHTTPResponse(ctx, rw, buf)
		return
	}
	params := r.FormValue("params")
	log.DEBUGF("params:%s", params)
	buf, _ := FormatResponse("Echo", 0, "", params)
	server.SendHTTPResponse(ctx, rw, buf)
	return
}
