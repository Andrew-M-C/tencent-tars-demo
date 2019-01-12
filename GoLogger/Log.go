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
	// "reflect"
	// "fmt"
	"strings"
	"path/filepath"
	"bytes"
	"strconv"
	"time"
	"sync"
	"github.com/TarsCloud/TarsGo/tars"
	rogger "github.com/TarsCloud/TarsGo/tars/util/rogger"
)

var log *rogger.Logger
var once sync.Once
var svcStr string
// var log = tars.GetLogger("service")
// var log = tars.GetRemoteLogger("service")

/**
 * reference: [Go 如何让函数只能被调用一次](https://blog.csdn.net/qq_36431213/article/details/83277009)
 */
func initLog() {
	cfg := tars.GetServerConfig()
	log = tars.GetRemoteLogger("service")
	svcStr = cfg.App + "." + cfg.Server
	log.Info("Log init")
	return
}

func packageLogText(level string, file string, line int, function string, text string) string {
	once.Do(initLog)
	now := time.Now()
	// microsecs := int((now.UnixNano() - now.Unix() * 1e9) / 1e3)
	time_str := now.Local().Format("2006-01-02 15:04:05.000000")

	var buffer bytes.Buffer
	buffer.WriteString(time_str)
	buffer.WriteString(" | ")
	buffer.WriteString(file)
	buffer.WriteString(", Line ")
	buffer.WriteString(strconv.Itoa(line))
	buffer.WriteString(" | ")
	buffer.WriteString(function)
	buffer.WriteString("() | ")
	buffer.WriteString(level)
	buffer.WriteString(" | ")
	buffer.WriteString(svcStr)
	buffer.WriteString(" | ")
	buffer.WriteString(text)

	// append tailing return character
	if ('\n' != text[len(text) - 1]) {
		buffer.WriteString("\n")
	}
	// fmt.Print(buffer.String())
	return buffer.String()
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
	file, line, function := getCallerInfo(0);
	log_msg := packageLogText("DEBUG", file, line, function, text)
	log.Debug(log_msg)
	// fmt.Println(log_msg)
}

func Info(text string) {
	file, line, function := getCallerInfo(0);
	log_msg := packageLogText("INFO ", file, line, function, text)
	log.Info(log_msg)
}

func Error(text string) {
	file, line, function := getCallerInfo(0);
	log_msg := packageLogText("ERROR", file, line, function, text)
	log.Error(log_msg)
}
