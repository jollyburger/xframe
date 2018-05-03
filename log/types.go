package log

const (
	defaultLogDir         = "log"
	defaultLogPrefix      = "xframe_"
	defaultLogSuffix      = ".log"
	defaultLogSize        = 50 // MB
	defaultLogLevelString = "DEBUG"
	defaultLogType        = "stdout"
)

const (
	Ldate         = 1 << iota                                 // the date: 2009/0123
	Ltime                                                     // the time: 01:23:23
	Lmicroseconds                                             // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                                                 // full file name and line number: /a/b/c/d.go:23
	Lshortfile                                                // final file name element and line number: d.go:23. overrides Llongfile
	Lmodule                                                   // module name
	Llevel                                                    // level: 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal)
	LstdFlags     = Ldate | Ltime | Lmicroseconds             // initial values for the standard logger
	Ldefault      = Lmodule | Llevel | Lshortfile | LstdFlags // [prefix][time][level][module][shortfile|longfile]
)

const (
	Lnop = iota
	Ldebug
	Linfo
	Lwarn
	Lerror
	Lpanic
	Lfatal
)

const BUFFER_SIZE = 1000

var levels = []string{
	"",
	"[DEBUG]",
	"[INFO]",
	"[WARN]",
	"[ERROR]",
	"[PANIC]",
	"[FATAL]",
}

var level_flags = []string{
	"",
	"debug",
	"info",
	"warn",
	"error",
	"panic",
	"fatal",
}

const (
	STACK_SIZE        = 64
	XFRAME_LOG_PREFIX = "xframe/log"
)
