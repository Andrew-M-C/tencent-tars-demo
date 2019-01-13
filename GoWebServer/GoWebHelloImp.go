package main

import (
	"fmt"
	"strings"
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/TarsCloud/TarsGo/tars"
	amc "amc/GoTarsServer/Amc"
	//log "github.com/Andrew-M-C/tencent-tars-demo/GoLogger"
	log "amc/GoLogger"
)

var comm *tars.Communicator

/**
 * HTTP return data
 * reference:
 * - [复合类型JSON / GO TYPES JSON - Author 品茶](https://www.kancloud.cn/zwhset/golang/363567)
 * - [golang的http cookie用法](http://www.cnblogs.com/liangDream/p/5522516.html)
 * - [golang web开发获取get、post、cookie参数](http://ju.outofmemory.cn/entry/93164)
 */
type HttpResp struct {
	Msg			string	`json:"msg"`
	Timestamp	int64	`json:"unix,omitempty"`
	TimeStr 	string 	`json:"time,omitempty"`
	Code		int		`json:"code,omitempty"`
	Client		string	`json:"client,omitempty"`
	User		string	`json:"user,omitempty"`
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

func HttpHelloHandler(w http.ResponseWriter, info *HttpRequestInfo, r *http.Request) {
	remote_ip, remote_port := getAddrFromRequest(r)
	log.Info(fmt.Sprintf("[%s:%d] remote http request", remote_ip, remote_port))

	// tars RPC
	comm = tars.NewCommunicator()
	app := new(amc.DateTime)
	obj := "amc.GoTarsServer.GoTarsObj"
	comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 10.0.4.11 -p 17890")

	req := amc.GetTimeReq{}
	rsp := amc.GetTimeRsp{}
	req.TimeFmt = "YYYY-MM-DD hh:mm:ss"
	var http_resp = HttpResp{}

	comm.StringToProxy(obj, app)
	ret, err := app.GetTime(&req, &rsp)
	if err != nil {
		log.Error(fmt.Sprintf("[%s:%d] Error, msg: %s, ret: %d", remote_ip, remote_port, err.Error(), ret))
		http_resp.Msg = err.Error()
		http_resp.Code = int(ret)
	} else {
		log.Debug(fmt.Sprintf("[%s:%d] Success, time %s", remote_ip, remote_port, rsp.LocalTimeStr))
		http_resp.Msg = "Hello, Tars-Go!"
		http_resp.Timestamp = int64(rsp.UtcTimestamp)
		http_resp.TimeStr = rsp.LocalTimeStr
	}

	// Other return messages
	http_resp.Client = fmt.Sprintf("%s:%d", remote_ip, remote_port)
	user_cookie, err := r.Cookie("user")
	if nil == err {
		user := strings.ToLower(user_cookie.Value)
		log.Debug("Cookie user: " + user);
		if user != "anonymous" {
			http_resp.User = user_cookie.Value
		}
	} else {
		log.Debug("Create a default user cookie");
		user_cookie = &http.Cookie{Name: "user", Value: "anonymous", HttpOnly: true, Secure: true, MaxAge: 0}
		http.SetCookie(w, user_cookie)
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	http_str, err := json.Marshal(http_resp)
	w.Write([]byte(http_str))
	return
}
