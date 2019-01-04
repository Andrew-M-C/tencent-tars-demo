package main

import (
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/TarsCloud/TarsGo/tars"
	amc "amc/GoTarsServer/Amc"
	_ "github.com/TarsCloud/TarsGo/tars/util/rogger"
)

var comm *tars.Communicator
var log = tars.GetLogger("log")

// HTTP return data
// reference: [复合类型JSON / GO TYPES JSON - Author 品茶](https://www.kancloud.cn/zwhset/golang/363567)
type HttpResp struct {
	Msg			string	`json:"msg"`
	Timestamp	int64	`json:"unix,omitempty"`
	TimeStr 	string 	`json:"time,omitempty"`
	Code		int		`json:"code,omitempty"`
	Client		string	`json:"client,omitempty"`
}

func getAddrFromRequest(r *http.Request) (ip string, port int) {
	if true {
		for key, value := range r.Header {
			log.Debug(fmt.Sprintf("[%s] - %s", key, value[0]))
		}
	}
	ip = r.Header.Get("X-Real-IP")
	port_str := r.Header.Get("X-Real-Port")
	port, err := strconv.Atoi(port_str)
	if nil != err {
		port = -1
	}
	return
}

func HttpRootHandler(w http.ResponseWriter, r *http.Request) {
	remote_ip, remote_port := getAddrFromRequest(r)
	log.Info(fmt.Sprintf("[%s:%d] remote http request", remote_ip, remote_port))

	comm = tars.NewCommunicator()
	app := new(amc.Go)
	obj := "amc.GoTarsServer.GoTarsObj"
	comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 10.0.4.11 -p 17890")

	req := amc.GetTimeReq{}
	rsp := amc.GetTimeRsp{}
	req.Time_fmt = "YYMMDD hh:mm:ss"
	var http_resp = HttpResp{}

	comm.StringToProxy(obj, app)
	ret, err := app.GetTime(&req, &rsp)
	if err != nil {
		log.Error(fmt.Sprintf("[%s:%d] Error, msg: %s, ret: %d", remote_ip, remote_port, err.Error(), ret))
		http_resp.Msg = err.Error()
		http_resp.Code = int(ret)
	} else {
		log.Debug(fmt.Sprintf("[%s:%d] Success, time %s", remote_ip, remote_port, rsp.Local_time_str))
		http_resp.Msg = "Hello, Tars-Go!"
		http_resp.Timestamp = int64(rsp.Utc_timestamp)
		http_resp.TimeStr = rsp.Local_time_str
	}
	http_resp.Client = fmt.Sprintf("%s:%d", remote_ip, remote_port)

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	http_str, err := json.Marshal(http_resp)
	w.Write([]byte(http_str))
	return
}
