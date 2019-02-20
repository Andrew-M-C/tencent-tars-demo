package main

import (
	"net/http"
	"github.com/Andrew-M-C/tencent-tars-demo/GoLogger/log"
	"github.com/TarsCloud/TarsGo/tars"
	"strings"
)

type HttpHandler func(http.ResponseWriter, *HttpRequestInfo, *http.Request)
var httpHandlers = map[string]HttpHandler {
	"/hello-tars":	HttpHelloHandler,
	"/tars":		HttpTarsHandler,
	"/tars/sendmsg":	HttpTarsSendMsgHandler,
}

/**
 * reference: [golang获取完整的url](https://www.tuicool.com/articles/juuu2qm)
 */
func httpRootHandler(w http.ResponseWriter, r *http.Request) {
	info := GetHttpRequestInfo(r)
	url := info.Url
	url_len := strings.Count(url, "") - 1
	if '/' == url[url_len - 1] && url_len > 1 {	// 删除结尾的 "/"
		url = url[0:url_len - 1]
		url_len -= 1
	}
	log.Infof("Request: %s, length %d", url, int(url_len))

	handler, found := httpHandlers[url]
	for false == found {
		// 截短 url 再进行搜索
		last_slash_index := strings.LastIndex(url, "/");
		if last_slash_index <= 0 {
			break;
		} else {
			url = url[0:last_slash_index];
			url_len = last_slash_index;
			log.Debug("Now check path: " + url)
			handler, found = httpHandlers[url]
		}
	}

	if false == found {
		log.Info(url + " not found")
		w.Write([]byte("404 " + url + " not found"))
	} else {
		handler(w, info, r)
	}
	log.Debug("== process ends ==")
}

func main() {
	cfg := tars.GetServerConfig()
	mux := &tars.TarsHttpMux{}
	mux.HandleFunc("/", httpRootHandler)
	tars.AddHttpServant(mux, cfg.App+"."+cfg.Server+".GoWebObj")

	tars.Run()
}
