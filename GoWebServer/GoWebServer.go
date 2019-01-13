package main

import (
	"net/http"
	log "amc/GoLogger"
	"github.com/TarsCloud/TarsGo/tars"
)

type HttpHandler func(http.ResponseWriter, *http.Request)
var httpHandlers = map[string]HttpHandler {
	"/hello-tars":	HttpHelloHandler,
	"/tars":		HttpTarsHandler,
}

/**
 * reference: [golang获取完整的url](https://www.tuicool.com/articles/juuu2qm)
 */
func httpRootHandler(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	url := r.RequestURI
	log.Info("Request: " + host + url);
	handler, found := httpHandlers[url]
	if false == found {
		w.Write([]byte("404 " + url + " not found"))
	} else {
		handler(w, r)
	}
}

func main() {
	cfg := tars.GetServerConfig()

	helloMux := &tars.TarsHttpMux{}
	helloMux.HandleFunc("/", httpRootHandler)
	tars.AddHttpServant(helloMux, cfg.App+"."+cfg.Server+".GoWebObj")

	tars.Run()
}
