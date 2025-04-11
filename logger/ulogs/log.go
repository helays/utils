package ulogs

import (
	"log"
	"os"
)

const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

var Level = 1

// Recover 捕获系统异常
func Recover() {
	if r := recover(); r != nil {
		Error("系统异常，捕获结果", r)
	}
}

// Log 打印正确日志，Info的别名
func Log(i ...interface{}) {
	Info(i...)
}

// Debug 用于记录调试信息
func Debug(i ...any) {
	if Level > LogLevelDebug {
		return
	}
	log.SetPrefix("【DEBUG】")
	log.SetOutput(os.Stdout)
	log.Println(i...)
}

func Debugf(format string, a ...any) {
	if Level > LogLevelDebug {
		return
	}
	log.SetPrefix("【DEBUG】")
	log.SetOutput(os.Stdout)
	log.Printf(format, a...)
}

// Info 用于记录信息
func Info(i ...interface{}) {
	if Level > LogLevelInfo {
		return
	}
	log.SetPrefix("【INFO】")
	log.SetOutput(os.Stdout)
	log.Println(i...)
}

func Infof(format string, a ...any) {
	if Level > LogLevelInfo {
		return
	}
	log.SetPrefix("【INFO】")
	log.SetOutput(os.Stdout)
	log.Printf(format, a...)
}

// Warn 用于记录警告信息
func Warn(i ...interface{}) {
	if Level > LogLevelWarn {
		return
	}
	log.SetPrefix("【WARN】")
	log.SetOutput(os.Stdout)
	log.Println(i...)
}

func Warnf(format string, a ...any) {
	if Level > LogLevelWarn {
		return
	}
	log.SetPrefix("【WARN】")
	log.SetOutput(os.Stdout)
	log.Printf(format, a...)
}

// Error 用于记录错误信息
func Error(i ...interface{}) {
	if Level > LogLevelError {
		return
	}
	log.SetPrefix("【ERROR】")
	log.SetOutput(os.Stderr)
	log.Println(i...)
}

func Errorf(format string, a ...any) {
	if Level > LogLevelError {
		return
	}
	log.SetPrefix("【ERROR】")
	log.SetOutput(os.Stderr)
	log.Printf(format, a...)
}

// Fatal 用于记录致命错误信息
func Fatal(i ...interface{}) {
	if Level > LogLevelFatal {
		return
	}
	log.SetPrefix("【FATAL】")
	log.SetOutput(os.Stderr)
	log.Println(i...)
}

func Fatalf(format string, a ...any) {
	if Level > LogLevelFatal {
		return
	}
	log.SetPrefix("【FATAL】")
	log.SetOutput(os.Stderr)
	log.Printf(format, a...)
}

// Checkerr 检查错误
func Checkerr(err error, i ...interface{}) {
	if err == nil {
		return
	}
	Error(append(i, err)...)
}

// DieCheckerr 检查错误，打印并输出错误信息
func DieCheckerr(err error, i ...any) {
	if err == nil {
		return
	}
	Error(append(i, err)...)
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
