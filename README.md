# xframe

xframe is a fast, scalable framework for golang. It provides basic development framework for different protocol and each components can be used separately.
xframe is suitable for api server, middleware and backend server.

---

## How to Install

Install & Run

```bash
go get github.com/jollyburger/xframe
```

---

## How to Use

Refer to example/http_example

in main.go, you need to initialize the config and start your socket service

```go
package main

import (
	"flag"
	"fmt"
	"xframe/cmd"
	"xframe/config"
	_ "xframe/example/http_example/logic"
	"xframe/handler/timer_handler"
	"xframe/log"
	"xframe/metric"
	"xframe/server"
	"xframe/trace"
)

var (
	confFile    = flag.String("c", "", "configuration file,json format")
	appName     = flag.String("a", "", "application name")
	confService = flag.String("s", "", "config service address,http server address")
)

func main() {
	//CMD Parser
	cmd.ParseCommand()
	cmd.DumpCommand()
	//Config initialization
	option, err := cmd.GetCommand("c")
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = config.LoadConfigFromFile(option); err != nil {
		fmt.Println("Load Config File fail,", err)
		return
	}
	config.DumpConfigContent()
	//Init to start pprof
	go metric.InitPprof("192.168.191.41:6060")
	//Init trace tool
	go trace.InitTraceTool()
	//Logger Init
	dir, _ := config.GetConfigByKey("log.LogDir")
	prefix, _ := config.GetConfigByKey("log.LogPrefix")
	suffix, _ := config.GetConfigByKey("log.LogSuffix")
	log_size, _ := config.GetConfigByKey("log.LogSize")
	log_level, _ := config.GetConfigByKey("log.LogLevel")
	log_type, _ := config.GetConfigByKey("log.LogType")
	log.InitLogger(dir.(string), prefix.(string), suffix.(string), int64(log_size.(float64)), log_level.(string), log_type.(string))
	// Timer start
	timer_handler.TimerTaskServe()
	// Service start
	ip, err := config.GetConfigByKey("http.listen_addr")
	if err != nil {
		fmt.Println("can not get listen ip:", err)
		return
	}
	listen_ip, _ := ip.(string)
	port, err := config.GetConfigByKey("http.listen_port")
	if err != nil {
		fmt.Println("can not get listen port")
		return
	}
	listen_port := int(port.(float64))
	if err = server.RunHTTP(listen_ip, listen_port); err != nil {
		fmt.Println(err)
		return
	}
}

```

for handler logic

```go
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
	//add customized data in ctx
	ctx := trace.InitContext()
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

```

---
## Tracing

### Build Jaeger 

refer to [jager-deployment](http://jaeger.readthedocs.io/en/latest/deployment/#configuration), building jaeger-agent, jaeger-collector, jaeger-query

refer to [jaeger-ui](https://github.com/uber/jaeger-ui), building jaeger-ui

### Add tracer 

#### Init Tracer
main.go

```
    import "xframe/trace"
    ... 
    
    trace.InitTracer("echo-server", "const", 1)
```

#### End to End
logic/init.go 

```
    import "xframe/trace"
    ...
    
    extracted_context, _ := trace.ExtractHTTPTracer(r.Header)
    sp := trace.StartSpan("tt", extracted_context)
    trace.FinishSpan(sp)
```

#### By Context
logic/echo.go

```
    import "xframe/trace"
    ...
    
    if span := trace.SpanFromContext(ctx); span != nil {
        span := trace.StartSpan("echo", span.Context())
        trace.SetSpanTag(span, "echo.content", "hello world")
        defer trace.FinishSpan(span)
        ctx = trace.ContextWithSpan(ctx, span)
    }   
```


---
## To be continued
