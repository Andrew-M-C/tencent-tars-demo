package main

import (
	"fmt"
	"net/http"
	"github.com/TarsCloud/TarsGo/tars"
	amc "amc/GoTarsServer/Amc"
)

var comm *tars.Communicator

func HttpRootHandler(w http.ResponseWriter, r *http.Request) {
	comm = tars.NewCommunicator()
	app := new(amc.Go)
	//obj := "amc.GoTarsServer.GoTarsObj@tcp -h 10.0.4.11 -p 10010 -t 60000"
	obj := "amc.GoTarsServer.GoTarsObj"
	comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 10.0.4.11 -p 17890")

	req := amc.GetTimeReq{}
	rsp := amc.GetTimeRsp{}
	resp_str := ""
	//resp_str := ""
	req.Time_fmt = "YYMMDD hh:mm:ss"

	comm.StringToProxy(obj, app)
	ret, err := app.GetTime(&req, &rsp)
	if err != nil {
		resp_str = fmt.Sprintf("{\"msg\":\"%s\", \"code\":%d}", err, ret)
	} else {
		resp_str = fmt.Sprintf("{\"msg\":\"Hello, Tars-Go!\", \"unix\":%d, \"time\":\"%s\"}", rsp.Utc_timestamp, rsp.Local_time_str)
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte(resp_str))
	return
}
