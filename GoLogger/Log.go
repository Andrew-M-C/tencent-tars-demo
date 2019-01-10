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
	"github.com/TarsCloud/TarsGo/tars"
)

var log = tars.GetLogger("service")

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
	var buffer bytes.Buffer
	buffer.WriteString(file)
	buffer.WriteString(", Line ")
	buffer.WriteString(strconv.Itoa(line))
	buffer.WriteString(", ")
	buffer.WriteString(function)
	buffer.WriteString("() - DEBUG - ")
	buffer.WriteString(text)
	log.Debug(buffer.String())
	// fmt.Println(buffer.String())
}

func Info(text string) {
	file, line, function := getCallerInfo(0);
	var buffer bytes.Buffer
	buffer.WriteString(file)
	buffer.WriteString(", Line ")
	buffer.WriteString(strconv.Itoa(line))
	buffer.WriteString(", ")
	buffer.WriteString(function)
	buffer.WriteString("() - INFO  - ")
	buffer.WriteString(text)
	log.Debug(buffer.String())
}

func Error(text string) {
	file, line, function := getCallerInfo(0);
	var buffer bytes.Buffer
	buffer.WriteString(file)
	buffer.WriteString(", Line ")
	buffer.WriteString(strconv.Itoa(line))
	buffer.WriteString(", ")
	buffer.WriteString(function)
	buffer.WriteString("() - ERROR - ")
	buffer.WriteString(text)
	log.Debug(buffer.String())
}
