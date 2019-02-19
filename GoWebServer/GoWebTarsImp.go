package main

import (
	// "fmt"
	// "strings"
	"net/http"
	// "encoding/json"
	// "github.com/TarsCloud/TarsGo/tars"
	//log "github.com/Andrew-M-C/tencent-tars-demo/GoLogger"
	log "amc/GoLogger"
)

func HttpTarsHandler(w http.ResponseWriter, info *HttpRequestInfo, r *http.Request) {
	log.Debug("MARK")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte("Hello, TarsGo"))
	return
}
