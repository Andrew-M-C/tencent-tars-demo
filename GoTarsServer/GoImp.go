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
	Amc "amc/GoTarsServer/Amc"
)

type GoImp struct {}
var log = tars.GetLogger("logic")

func (imp *GoImp) GetTime(req *Amc.GetTimeReq, rsp *Amc.GetTimeRsp) (int32, error) {
	utc_time := time.Now()
	local_time := utc_time.Local()

	// convert time string
	var time_str string
	if "" == (*req).Time_fmt {
		log.Info("Use default time format")
		time_str = local_time.Format("01/02 15:04:05 2006")
	} else {
		log.Info(fmt.Sprintf("Got format string: %s", (*req).Time_fmt))
		time_str = (*req).Time_fmt
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

	(*rsp).Utc_timestamp = utc_time.Unix()
	(*rsp).Local_timestamp = local_time.Unix()
	(*rsp).Local_time_str = time_str
	return 0, nil
};
