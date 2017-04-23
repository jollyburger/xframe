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
