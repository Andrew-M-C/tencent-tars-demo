package main

import (
	"net/http"
	"github.com/Andrew-M-C/tencent-tars-demo/GoLogger/log"
	"github.com/TarsCloud/TarsGo/tars"
	"strings"
	"runtime"
	"github.com/capnm/sysinfo"
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
	server := cfg.App + "." + cfg.Server
	mux.HandleFunc("/", httpRootHandler)
	tars.AddHttpServant(mux, server + ".GoWebObj")

	{
		// ref: [golang获取服务内存信息](https://blog.csdn.net/m0_38132420/article/details/71699815)
		// ref: [package sysinfo](https://godoc.org/github.com/capnm/sysinfo#SI.ToString)
		sys := sysinfo.Get()

		secs := sys.Uptime / 1000000000
		day := secs / (60*60*24)
		secs -= day * (60*60*24)
		hour := secs / (60*60)
		secs -= hour * (60*60)
		min := secs / 60
		secs -= min * 60

		log.Infof("%s ready to run, CPU core(s): %d; sysup: %dd %dh %dm %ds; memory: %d MB",
					server, runtime.NumCPU(), day, hour, min, secs, sys.TotalRam >> 10)
	}
	tars.Run()
}
