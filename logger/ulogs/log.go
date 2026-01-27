package ulogs

import (
	"log"
	"os"
)

const (
	LogLevelTrace = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

var (
	traceLogger = log.New(os.Stdout, "【TRACE】", log.LstdFlags)
	debugLogger = log.New(os.Stdout, "【DEBUG】", log.LstdFlags)
	infoLogger  = log.New(os.Stdout, "【INFO】", log.LstdFlags)
	warnLogger  = log.New(os.Stdout, "【WARN】", log.LstdFlags)
	errorLogger = log.New(os.Stderr, "【ERROR", log.LstdFlags)
	fatalLogger = log.New(os.Stderr, "【FATAL】", log.LstdFlags)
)

var Level = LogLevelInfo

// Recover 捕获系统异常
func Recover() {
	if r := recover(); r != nil {
		Error("系统异常，捕获结果", r)
	}
}

// Log 打印正确日志，Info的别名
// Deprecated: 弃用,请使用 Info
func Log(i ...interface{}) {
	Info(i...)
}

func Trace(i ...any) {
	if Level <= LogLevelTrace {
		traceLogger.Println(i...)
	}
}

// noinspection all
func Tracef(format string, a ...any) {
	if Level <= LogLevelTrace {
		traceLogger.Printf(format, a...)
	}
}

// Debug 用于记录调试信息
func Debug(i ...any) {
	if Level <= LogLevelDebug {
		debugLogger.Println(i...)
	}
}

// Debugf
// noinspection all
func Debugf(format string, a ...any) {
	if Level <= LogLevelDebug {
		debugLogger.Printf(format, a...)
	}
}

// Info 用于记录信息
func Info(i ...interface{}) {
	if Level <= LogLevelInfo {
		infoLogger.Println(i...)
	}
}

// noinspection all
func Infof(format string, a ...any) {
	if Level <= LogLevelInfo {
		infoLogger.Printf(format, a...)
	}

}

// Warn 用于记录警告信息
func Warn(i ...interface{}) {
	if Level <= LogLevelWarn {
		warnLogger.Println(i...)
	}
}

// noinspection all
func Warnf(format string, a ...any) {
	if Level <= LogLevelWarn {
		warnLogger.Printf(format, a...)
	}
}

// Error 用于记录错误信息
func Error(i ...interface{}) {
	if Level <= LogLevelError {
		errorLogger.Println(i...)
	}
}

func Errorf(format string, a ...any) {
	if Level <= LogLevelError {
		errorLogger.Printf(format, a...)
	}
}

// Fatal 用于记录致命错误信息
func Fatal(i ...interface{}) {
	if Level <= LogLevelFatal {
		fatalLogger.Println(i...)
	}
}

// noinspection all
func Fatalf(format string, a ...any) {
	if Level <= LogLevelFatal {
		fatalLogger.Printf(format, a...)
	}
}

// Checkerr 检查错误
func Checkerr(err error, i ...interface{}) {
	if err == nil {
		return
	}
	Error(append(i, err)...)
}

// CheckErrf 检查错误
func CheckErrf(err error, format string, a ...any) {
	if err == nil {
		return
	}
	Errorf(format+" 原始错误 %v", append(a, err)...)
}

// DieCheckerr 检查错误，打印并输出错误信息
func DieCheckerr(err error, i ...any) {
	if err == nil {
		return
	}
	Error(append(i, err)...)
	os.Exit(1)
}

func DieCheckErrf(err error, format string, a ...any) {
	if err == nil {
		return
	}
	Errorf(format+" 原始错误 %v", append(a, err)...)
	os.Exit(1)
}

// ReturnCheckerr 检查错误，有异常就返回false
func ReturnCheckerr(err error, i ...interface{}) bool {
	if err == nil {
		return true
	}
	Error(append(i, err)...)
	return false
}

func ErrorReturn(i ...interface{}) bool {
	Error(i...)
	return false
}

func Pfunc(a ...interface{}) {
	// log.SetPrefix("[用户异常]")
	log.SetOutput(os.Stdout)
	log.Println(a...)
}
