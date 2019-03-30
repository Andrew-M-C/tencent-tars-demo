package main

import (
	"github.com/Andrew-M-C/go-tools/log"
	"net/http"
	"utils/singleproc"
	httpstatic "server/static"
	httpdynamic "server/dynamic"
)

// func httpLoginHandler(w http.ResponseWriter, r *http.Request) {
// 	info := httpparser.GetHttpRequestInfo(r)
// 	url := info.Url
// 	log.Debug("request URL: %s", url)
// 	useful_url := strings.Replace(url, "/", "", 1)
// 	log.Debug("request resource: %s", useful_url)

// 	for k, v := range info.Query {
// 		log.Debug("%s = %s", k, v)
// 	}

// 	w.Header().Set("Content-Type", "application/json;charset=utf-8")
// 	w.Write([]byte(`{"ret":0, "message":"OK"}`))
// 	return
// }


func main() {
	pid, err := singleproc.GetRunningProcessID()
	if err != nil {
		log.Error("Failed tp get pid: %s", err.Error())
	} else if pid > 0 {
		log.Info("Rnning process pid: %d, we should quit", pid)
		return
	}

	// run in backgound
	pid, err = singleproc.DaemonizeAndLogPid()
	if err != nil {
		log.Error("Failed to daemonize: %s", err.Error())
		return
	} else {
		log.Info("Daemonize successfully, pid: %d", pid)
	}

	http.HandleFunc("/", httpstatic.RootHandler)
	http.HandleFunc("/html/", httpstatic.HtmlHandler)
	http.HandleFunc("/js/", httpstatic.JsHandler)
	http.HandleFunc("/cgi-bin/", httpdynamic.LoginHandler)
	http.HandleFunc("/html/success.html", httpdynamic.SuccessHandler)
	http.ListenAndServe(":4000", nil)
	log.Info("process %d quit", pid)
	return
}
