package logger

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

type LogLevel uint16

const (
	UNKNOWN LogLevel = iota
	DEBUG
	INFO
	WARNING
	ERROR
)

// 将字符串转为对应日志级别LogLevel
func parseLogLevel(s string) (LogLevel, error) {
	switch strings.ToLower(s) {
	case "debug":
		return DEBUG, nil
	case "info":
		return INFO, nil
	case "warning":
		return WARNING, nil
	case "error":
		return ERROR, nil
	default:
		err := errors.New("invalid log level")
		return UNKNOWN, err
	}
}

// 将日志级别LogLevel转为对应字符串
func getLogString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	default:
		return "DEBUG"
	}
}

// 获取当前时间,Y-m-d H:i:s
func getTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 获取当前执行位置信息
// skip = 0 ,getInfo() log_const.go 63
// skip = 1 ,(*FileLogger).logAsync() log_handle.go 178
// skip = 2 ,(*FileLogger).Info() log_const.go 257
// skip = 3 ,TestLogger() logger_test.go 9
func getInfo(skip int) (funcName, fileName string, lineNo int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("runtime.Caller() failed")
		return
	}
	return strings.Split(runtime.FuncForPC(pc).Name(), ".")[1], path.Base(file), line
}
