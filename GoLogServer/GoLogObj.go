/**
 * reference:
 * - [golang time的几种用法](https://my.oschina.net/yinlei212/blog/151963)
 */
package main

import (
	"sync"
	"strings"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/util/rogger"
	log "./Logf"
	logClient "../GoLogger/log"
)

type GoLogObj struct {}
var once sync.Once
var llog *rogger.Logger

func init() {
	llog = tars.GetLogger("remote_logs")
	return
}

func (imp *GoLogObj) LoggerbyInfo(logConf *log.LogInfo, logList []string) error {
	for _, each_log := range logList {
		seperator := strings.Index(each_log, "|") + 1
		log_msg, _ := logClient.ParseLogText(each_log[seperator:])
		if strings.Count(log_msg.Datetime, "") > 0 {
			llog.Infof(" - %s.%s <%s> | %s | %s:%d | %s() | %s | %s",
						logConf.Appname, logConf.Servername, logConf.SFilename,
						log_msg.Datetime, log_msg.File, int(log_msg.Line), log_msg.Function,
						log_msg.LevelStr, log_msg.Text)
		}
	}
	return nil
};

func (imp *GoLogObj) Logger(app string, server string, file string, format string, buffer []string) error {
	llog.Info("MARK")
	return nil
}
