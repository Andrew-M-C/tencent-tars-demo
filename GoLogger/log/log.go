/**
 * reference:
 * - [go语言读取当前文件名行号和函数名](https://www.jianshu.com/p/7aa54489bb4b)
 * - [如何在Go中获取函数的名称？](https://codeday.me/bug/20170831/60552.html)
 * - [golang 几种字符串的拼接方式](https://blog.csdn.net/iamlihongwei/article/details/79551259)
 * - [基础知识 - Golang 中的格式化输入输出](https://www.cnblogs.com/golove/p/3284304.html)
 */
package log

import (
	"runtime"
	"strings"
	"path/filepath"
	"encoding/json"
	"time"
	"fmt"
	// "sync"
	"github.com/TarsCloud/TarsGo/tars"
	rogger "github.com/TarsCloud/TarsGo/tars/util/rogger"
)

var log *rogger.Logger
var shouldLogDebug = false
var shouldLogInfo = true
var shouldLogWarn = true
var shouldLogError = true
// var once sync.Once
// var svcStr string
// var log = tars.GetLogger("service")
// var log = tars.GetRemoteLogger("service")

/**
 * reference: [Go 如何让函数只能被调用一次](https://blog.csdn.net/qq_36431213/article/details/83277009)
 */
func init() {
	log = tars.GetRemoteLogger("service")
	cfg := tars.GetServerConfig()
	// svcStr = cfg.App + "." + cfg.Server
	// determine log level
	log_level := cfg.LogLevel
	if log_level == "DEBUG" {
		shouldLogDebug = true
	} else if log_level == "INFO" {
		// nothing to change
	} else if log_level == "WARN" || log_level == "WARNING" {
		shouldLogInfo = false
	} else if log_level == "ERROR" || log_level == "ERR" {
		shouldLogInfo = false
		shouldLogWarn = false
	} else if log_level == "OFF" {
		shouldLogInfo = false
		shouldLogWarn = false
		shouldLogError = false
	} else {
		// nothing to do
	}
	log.Info("Log level: " + log_level)
	return
}

func packageLogText(level int32, file string, line int, function string, text string) string {
	var log_text LogText
	now := time.Now()
	log_text.Datetime = now.Local().Format("2006-01-02 15:04:05.000000")
	log_text.Timestamp = float64(now.UnixNano() / 1e3) / 1e6
	log_text.LevelInt = level
	log_text.LevelStr = GetLevelString(level)
	log_text.File = file
	log_text.Line = int32(line)
	log_text.Function = function
	log_text.Text = text

	json_bytes, _ := json.Marshal(log_text)
	return string(json_bytes)
}

func getCallerInfo(invoke_level int) (fileName string, line int, funcName string) {
	funcName = "FILE"
	line = -1
	fileName = "FUNC"

	if invoke_level <= 0 {
		invoke_level = 2
	} else {
		invoke_level += 1
	}

	pc, file_name, line, ok := runtime.Caller(invoke_level)
	if ok {
		fileName = filepath.Base(file_name)
		func_name := runtime.FuncForPC(pc).Name()
		func_name = filepath.Ext(func_name)
		funcName = strings.TrimPrefix(func_name, ".")
	}
	// fmt.Println(reflect.TypeOf(pc), reflect.ValueOf(pc))
	return
}

func Debug(text string) {
	if false == shouldLogDebug {
		return
	}
	file, line, function := getCallerInfo(0);
	log_msg := packageLogText(LEVEL_DEBUG, file, line, function, text)
	log.Debug(log_msg)
	// fmt.Println(log_msg)
}

func Info(text string) {
	if false == shouldLogInfo {
		return
	}
	file, line, function := getCallerInfo(0);
	log_msg := packageLogText(LEVEL_INFO, file, line, function, text)
	log.Info(log_msg)
}

func Warn(text string) {
	if false == shouldLogWarn {
		return
	}
	file, line, function := getCallerInfo(0);
	log_msg := packageLogText(LEVEL_WARN, file, line, function, text)
	log.Error(log_msg)
}

func Error(text string) {
	if false == shouldLogError {
		return
	}
	file, line, function := getCallerInfo(0);
	log_msg := packageLogText(LEVEL_ERROR, file, line, function, text)
	log.Error(log_msg)
}

func Debugf(format string, v ...interface{}) {
	if false == shouldLogDebug {
		return
	}
	file, line, function := getCallerInfo(0);
	text := fmt.Sprintf(format, v...)
	log_msg := packageLogText(LEVEL_DEBUG, file, line, function, text)
	log.Debug(log_msg)
}

func Infof(format string, v ...interface{}) {
	if false == shouldLogInfo {
		return
	}
	file, line, function := getCallerInfo(0);
	text := fmt.Sprintf(format, v...)
	log_msg := packageLogText(LEVEL_INFO, file, line, function, text)
	log.Info(log_msg)
}

func Warnf(format string, v ...interface{}) {
	if false == shouldLogWarn {
		return
	}
	file, line, function := getCallerInfo(0);
	text := fmt.Sprintf(format, v...)
	log_msg := packageLogText(LEVEL_WARN, file, line, function, text)
	log.Warn(log_msg)
}

func Errorf(format string, v ...interface{}) {
	if false == shouldLogInfo {
		return
	}
	file, line, function := getCallerInfo(0);
	text := fmt.Sprintf(format, v...)
	log_msg := packageLogText(LEVEL_ERROR, file, line, function, text)
	log.Error(log_msg)
}
