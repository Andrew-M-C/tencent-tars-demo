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
	_ "reflect"
	_ "fmt"
	"strings"
	"path/filepath"
	"bytes"
	"strconv"
	"time"
	"github.com/TarsCloud/TarsGo/tars"
)

var log = tars.GetLogger("service")

func packageLogText(level string, file string, line int, function string, text string) string {
	now := time.Now()
	milisecs := int((now.UnixNano() - now.Unix() * 1e9) / 1e6)
	time_str := now.Local().Format("2006-01-02 15:04:05.")

	var buffer bytes.Buffer
	buffer.WriteString(time_str)
	if (milisecs >= 100) {
		buffer.WriteString(strconv.Itoa(milisecs))
	} else if (milisecs >= 10) {
		buffer.WriteString("0")
		buffer.WriteString(strconv.Itoa(milisecs))
	} else {
		buffer.WriteString("00")
		buffer.WriteString(strconv.Itoa(milisecs))
	}

	buffer.WriteString(" | ")
	buffer.WriteString(file)
	buffer.WriteString(" | Line ")
	buffer.WriteString(strconv.Itoa(line))
	buffer.WriteString(" | ")
	buffer.WriteString(function)
	buffer.WriteString("() | ")
	buffer.WriteString(level)
	buffer.WriteString(" | ")
	buffer.WriteString(text)
	return buffer.String()
}

func getCallerInfo(invoke_level int) (fileName string, line int, funcName string) {
	funcName = "unknown_func"
	line = -1
	fileName = "unknown_file"
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
