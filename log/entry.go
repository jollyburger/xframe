package log

import (
	"fmt"
	"os"
)

var (
	Glogger *Logger
)

func InitLogger(dir, prefix, suffix string, size int64, level string, logtype string) {
	if Glogger != nil {
		Glogger.Close()
	}
	if logtype == "" {
		logtype = defaultLogType
	}
	if level == "" {
		level = defaultLogLevelString
	}
	var (
		logger *Logger
		err    error
	)
	switch logtype {
	case "file":
		if dir == "" {
			dir = defaultLogDir
		}
		if prefix == "" {
			prefix = defaultLogPrefix
		}
		if suffix == "" {
			suffix = defaultLogSuffix
		}
		if size <= 0 {
			size = defaultLogSize
		}
		logger, err = NewRotate(dir, prefix, suffix, size)
		if err != nil {
			fmt.Println("Init Logger fail:", err)
			os.Exit(-1)
		}
	case "stdout":
		logger = New(os.Stdout, "", Ldefault)
	}
	Glogger = logger
	SetLogLevel(level)
}

func InitDefaultLogger() {
	InitLogger(defaultLogDir, defaultLogPrefix, defaultLogSuffix, defaultLogSize, defaultLogLevelString, "")
}

func SetLogLevel(level string) {
	if Glogger == nil {
		InitLogger(defaultLogDir, defaultLogPrefix, defaultLogSuffix, defaultLogSize, level, "")
	}
	Glogger.SetOutputLevelString(level)
}

func SetDailyRotate(daily bool) {
	if Glogger == nil {
		InitLogger(defaultLogDir, defaultLogPrefix, defaultLogSuffix, defaultLogSize, defaultLogLevelString, "")
	}
	Glogger.SetDailyRotate(daily)
}

func SetBackup(backup int) {
	if Glogger == nil {
		InitLogger(defaultLogDir, defaultLogPrefix, defaultLogSuffix, defaultLogSize, defaultLogLevelString, "")
	}
	Glogger.SetBackup(backup)
}

func EnableLogDepth(flag bool) {
	if Glogger == nil {
		InitLogger(defaultLogDir, defaultLogPrefix, defaultLogSuffix, defaultLogSize, defaultLogLevelString, "")
	}
	Glogger.enableLogDepth(flag)
}

func SetHook(hook Hook) error {
	if Glogger == nil {
		InitLogger(defaultLogDir, defaultLogPrefix, defaultLogSuffix, defaultLogSize, defaultLogLevelString, "")
	}
	return Glogger.rotateLogger.AddHook(hook)
}

//====================================================

func INFOF(format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Infof(format, v...)
}

func INFO(v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Info(v...)
}

func ERRORF(format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Errorf(format, v...)
}

func ERROR(v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Error(v...)
}

func WARN(v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Warn(v...)
}

func WARNF(format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Warnf(format, v...)
}

func DEBUG(v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Debug(v...)
}

func DEBUGF(format string, v ...interface{}) {
	if Glogger == nil {
		InitDefaultLogger()
	}
	Glogger.Debugf(format, v...)
}
