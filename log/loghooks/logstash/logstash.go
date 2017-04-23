package uflog_logstash

import (
	"strconv"
	"strings"
	"xframe/log/loghooks/logstash/go-logstash"
)

const (
	LOGSTASH_TIMEOUT = 10
)

type LogstashHook struct {
	LogstashInst *logstash.Logstash
	LevelLst     []string
	Addr         string
}

func NewLogstashHook(addr string, levels []string) (logstash_hook *LogstashHook, err error) {
	hostname := strings.Split(addr, ":")[0]
	port, _ := strconv.Atoi(strings.Split(addr, ":")[1])
	logstash_instance := logstash.New(hostname, port, LOGSTASH_TIMEOUT)
	logstash_hook = &LogstashHook{
		LogstashInst: logstash_instance,
		LevelLst:     levels,
		Addr:         addr,
	}
	_, err = logstash_instance.Connect()
	if err != nil {
		return
	}
	return
}

func (lsh LogstashHook) Levels() []string {
	return lsh.LevelLst
}

func (lsh LogstashHook) Fire(level string, msg []byte) error {
	var err error
	err = lsh.LogstashInst.Writeln(string(msg))
	if err != nil {
		//reconnect and retry once
		_, err = lsh.LogstashInst.Connect()
		if err != nil {
			goto END
		}
		err = lsh.LogstashInst.Writeln(string(msg))
		if err != nil {
			goto END
		}
	}
END:
	return err
}
