/**
 * reference:
 * - [golang time的几种用法](https://my.oschina.net/yinlei212/blog/151963)
 */
package main

import (
	_ "fmt"
	_ "time"
	_ "strings"
	"sync"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/util/rogger"
	log "./Logf"
)

type GoLogObj struct {}
var once sync.Once
var llog *rogger.Logger

func initLog() {
	llog = tars.GetLogger("remote_logs")
	return
}

func (imp *GoLogObj) LoggerbyInfo(logConf *log.LogInfo, logList []string) error {
	once.Do(initLog)
	for _, each_log := range logList {
		llog.Info(logConf.Appname + "." + logConf.Servername + ", " + logConf.SFilename + ": " + each_log)
	}
	return nil
};

func (imp *GoLogObj) Logger(app string, server string, file string, format string, buffer []string) error {
	once.Do(initLog)
	llog.Info("MARK")
	return nil
}
