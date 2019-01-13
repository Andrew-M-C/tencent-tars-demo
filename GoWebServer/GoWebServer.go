package main

import (
	"net/http"
	log "amc/GoLogger"
	"github.com/TarsCloud/TarsGo/tars"
)

type HttpHandler func(http.ResponseWriter, *HttpRequestInfo, *http.Request)
var httpHandlers = map[string]HttpHandler {
	"/hello-tars":	HttpHelloHandler,
	"/tars":		HttpTarsHandler,
}

/**
 * reference: [golang获取完整的url](https://www.tuicool.com/articles/juuu2qm)
 */
func httpRootHandler(w http.ResponseWriter, r *http.Request) {
	info := GetHttpRequestInfo(r)
	log.Info("Request: " + info.Url);
	handler, found := httpHandlers[info.Url]
	if false == found {
		w.Write([]byte("404 " + info.Url + " not found"))
	} else {
		handler(w, info, r)
	}
}

func main() {
	cfg := tars.GetServerConfig()
	mux := &tars.TarsHttpMux{}
	mux.HandleFunc("/", httpRootHandler)
	tars.AddHttpServant(mux, cfg.App+"."+cfg.Server+".GoWebObj")

	tars.Run()
}
