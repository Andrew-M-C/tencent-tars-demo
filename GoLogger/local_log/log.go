/**
 */
package log

import (
	"runtime"
	"strings"
	"path/filepath"
	"time"
	"fmt"
)

func packageLogText(level string, file string, line int, function string, text string) string {
	now := time.Now()
	time_str := now.Local().Format("2006-01-02 15:04:05.000000")
	return fmt.Sprintf("%s - %s | %s, Line %d | %s() | %s", time_str, level, file, line, function, text)
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
	fmt.Println(log_msg)
}

func Info(text string) {
	file, line, function := getCallerInfo(0);
	log_msg := packageLogText("INFO ", file, line, function, text)
	fmt.Println(log_msg)
}

func Warn(text string) {
	file, line, function := getCallerInfo(0);
	log_msg := packageLogText("WARN ", file, line, function, text)
	fmt.Println(log_msg)
}

func Error(text string) {
	file, line, function := getCallerInfo(0);
	log_msg := packageLogText("ERROR", file, line, function, text)
	fmt.Println(log_msg)
}

func Debugf(format string, v ...interface{}) {
	file, line, function := getCallerInfo(0);
	text := fmt.Sprintf(format, v...)
	log_msg := packageLogText("DEBUG", file, line, function, text)
	fmt.Println(log_msg)
}

func Infof(format string, v ...interface{}) {
	file, line, function := getCallerInfo(0);
	text := fmt.Sprintf(format, v...)
	log_msg := packageLogText("INFO ", file, line, function, text)
	fmt.Println(log_msg)
}

func Warnf(format string, v ...interface{}) {
	file, line, function := getCallerInfo(0);
	text := fmt.Sprintf(format, v...)
	log_msg := packageLogText("WARN ", file, line, function, text)
	fmt.Println(log_msg)
}

func Errorf(format string, v ...interface{}) {
	file, line, function := getCallerInfo(0);
	text := fmt.Sprintf(format, v...)
	log_msg := packageLogText("ERROR", file, line, function, text)
	fmt.Println(log_msg)
}
