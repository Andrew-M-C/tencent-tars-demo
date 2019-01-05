/**
 * reference:
 * - [golang time的几种用法](https://my.oschina.net/yinlei212/blog/151963)
 */
package main

import (
	"fmt"
	"time"
	"strings"
	"github.com/TarsCloud/TarsGo/tars"
	amc "amc/GoTarsServer/Amc"
)

type GoImp struct {}
var log = tars.GetLogger("logic")

func (imp *GoImp) GetTime(req *amc.GetTimeReq, rsp *amc.GetTimeRsp) (int32, error) {
	utc_time := time.Now()
	local_time := utc_time.Local()

	// convert time string
	var time_str string
	if "" == (*req).TimeFmt {
		log.Info("Use default time format")
		time_str = local_time.Format("01/02 15:04:05 2006")
	} else {
		/**
		 * reference:
		 * - [go 时间格式风格详解](https://my.oschina.net/achun/blog/142315)
		 * - [Go 时间格式化和解析](https://www.kancloud.cn/itfanr/go-by-example/81698)
		 */
		log.Info(fmt.Sprintf("Got format string: %s", (*req).TimeFmt))
		time_str = (*req).TimeFmt
		time_str = strings.Replace(time_str, "YYYY", "2006", -1)
		time_str = strings.Replace(time_str, "yyyy", "2006", -1)
		time_str = strings.Replace(time_str, "YY", "06", -1)
		time_str = strings.Replace(time_str, "yy", "06", -1)
		time_str = strings.Replace(time_str, "MM", "01", -1)
		time_str = strings.Replace(time_str, "dd", "02", -1)
		time_str = strings.Replace(time_str, "DD", "02", -1)
		time_str = strings.Replace(time_str, "hh", "15", -1)
		time_str = strings.Replace(time_str, "mm", "04", -1)
		time_str = strings.Replace(time_str, "ss", "05", -1)
		log.Info("Convert as golang format: ", time_str)
		time_str = local_time.Format(time_str)
	}

	(*rsp).UtcTimestamp = utc_time.Unix()
	(*rsp).LocalTimestamp = local_time.Unix()
	(*rsp).LocalTimeStr = time_str
	return 0, nil
};
