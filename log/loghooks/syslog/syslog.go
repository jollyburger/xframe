package ulog_syslog

import (
	"errors"
	"log/syslog"
)

type SyslogHook struct {
	Writer *syslog.Writer
	Levels []string
	Addr   string
}

func SyslogPrio(level string) syslog.Priority {
	switch level {
	case "debug":
		return syslog.LOG_DEBUG
	case "info":
		return syslog.LOG_INFO
	case "error":
		return syslog.LOG_ERR
	case "warn":
		return syslog.LOG_WARNING
	default:
		return syslog.LOG_DEBUG
	}
}

func NewSyslogHook(addr string, levels []string) (syslog_hook *SyslogHook, err error) {
	writer, err := syslog.Dial("tcp", addr, syslog.LOG_DEBUG, "")
	if err != nil {
		return
	}
	syslog_hook = &SyslogHook{
		Writer: writer,
		Levels: levels,
		Addr:   addr,
	}
	return
}

func (slh SyslogHook) GetLevels() []string {
	return slh.Levels
}

func (slh SyslogHook) Fire(level string, msg []byte) error {
	switch level {
	case "debug":
		return slh.Writer.Debug(string(msg))
	case "info":
		return slh.Writer.Info(string(msg))
	case "error":
		return slh.Writer.Err(string(msg))
	case "warn":
		return slh.Writer.Warning(string(msg))
	default:
		return errors.New("not set hook level")
	}
}
