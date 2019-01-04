/**
 * reference:
 * - [golang time的几种用法](https://my.oschina.net/yinlei212/blog/151963)
 */
package main

import (
	"time"
	Amc "amc/GoTarsServer/Amc"
)

type GoImp struct {}

func (imp *GoImp) GetTime(req *Amc.GetTimeReq, rsp *Amc.GetTimeRsp) (int32, error) {
	utc_time := time.Now()
	local_time := utc_time.Local()
	time_str := local_time.Format("2006-01-02 15:04:05")

	(*rsp).Utc_timestamp = utc_time.Unix()
	(*rsp).Local_timestamp = local_time.Unix()
	(*rsp).Local_time_str = time_str
	return 0, nil
};
