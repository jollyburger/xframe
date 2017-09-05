package main

import (
	"flag"
	"xframe/cmd"
	"xframe/config"
	"xframe/log"
	"xframe/server"
)

var (
	confFile = flag.String("c", "", "configuration file, json format")
)

func main() {
	cmd.ParseCommand()
	cmd.DumpCommand()
	option, err := cmd.GetCommand("c")
	if err != nil {
		log.ERROR(err)
		panic(err)
	}
	if err := config.LoadConfigFromFile(option); err != nil {
		log.ERROR(err)
		panic(err)
	}
	config.DumpConfigContent()

	addr, _ := config.GetConfigByKey("server.ip")
	port, _ := config.GetConfigByKey("server.port")

	mux := handler.RegisterHandler()
	if err := server.RunHTTPMux(addr.(string), int(port.(float64)), mux); err != nil {
		log.ERROR(err)
		panic(err)
	}
}
