package main

import (
	"fmt"
	// "strings"
	"strconv"
	"net/http"
	// "encoding/json"
	// "github.com/TarsCloud/TarsGo/tars"
	//log "github.com/Andrew-M-C/tencent-tars-demo/GoLogger"
	log "amc/GoLogger"
)

type httpRequestInfo struct {
	Ip			string
	Port		int
	Url			string
	Header		map[string]string
	Query		map[string]string
}

func parseHttpRequest(r *http.Request, info *httpRequestInfo) {
	if true {
		for key, value := range r.Header {
			log.Debug(fmt.Sprintf("[%s] - %s", key, value[0]))
			info.Header[key] = value[0]
		}
	}
	info.Ip = r.Header.Get("X-Real-IP")
	port_str := r.Header.Get("X-Real-Port")
	port, err := strconv.Atoi(port_str)
	if nil != err {
		port = -1
	}
	info.Port = port
	return
}

func HttpTarsHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("MARK")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte("Hello, TarsGo"))
	return
}
