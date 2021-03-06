package main

import (
	"net/http"
	"github.com/Andrew-M-C/tencent-tars-demo/GoLogger/log"
	// log "../GoLogger/local_log"
	"github.com/TarsCloud/TarsGo/tars"
	"strings"
	"runtime"
	"github.com/capnm/sysinfo"
	"github.com/Andrew-M-C/tarsgo-tools/config"
	// "encoding/json"
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
	// configurations
	cfg := tars.GetServerConfig()
	mux := &tars.TarsHttpMux{}
	server := cfg.App + "." + cfg.Server

	// check self-defined configurations
	{
		tarsconf, err := config.NewConfig()
		if err != nil {
			log.Debug("Failed to get config: " + err.Error())
		} else {
			log_path, exist := tarsconf.GetString("/tars/application/server", "logpath", "unknown")
			log.Debugf("%t, logpath: %s", exist, log_path)

			log_level, exist := tarsconf.GetString("/tars/application/server", "logLevel", "DEBUG")
			log.Debugf("%t, logLevel: %s", exist, log_level)
		}
	}

	// add servants
	{
		mux.HandleFunc("/", httpRootHandler)
		tars.AddHttpServant(mux, server + ".GoWebObj")
	}

	// check OS status
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

		log.Infof("%s ready to run, CPU core(s) | %d; sysup: %dd %dh %dm %ds | memory: %d MB",
					server, runtime.NumCPU(), day, hour, min, secs, sys.TotalRam >> 10)
	}

	// This is an example for how to parse a structure-unknown JSON string
	// additional ref: [preserve int64 values when parsing json in Go](https://stackoverflow.com/questions/16946306/preserve-int64-values-when-parsing-json-in-go)
	// {
	// 	var func_parse_array func([]interface{}, int)
	// 	var func_parse_obj func(map[string]interface{}, int)

	// 	func_parse_array = func(array []interface{} ,level int) {
	// 		prefix := strings.Repeat("    ", level)
	// 		for index, intf := range array {
	// 			switch intf.(type) {
	// 			case string:
	// 				log.Infof(prefix + "[%d] = %s", index, intf.(string))
	// 			case float64:
	// 				integer := int64(intf.(float64))
	// 				if (float64(integer) == intf.(float64)) {
	// 					log.Infof(prefix + "[%d] = %d", index, integer)
	// 				} else {
	// 					log.Infof(prefix + "[%d] = %f", index, intf.(float64))
	// 				}
	// 			case bool:
	// 				log.Infof(prefix + "[%d] = %t", index, intf.(bool))
	// 			case map[string]interface{}:
	// 				log.Infof(prefix + "[%d] is an object", index)
	// 				func_parse_obj(intf.(map[string]interface{}), level + 1)
	// 			case []interface{}:
	// 				log.Infof(prefix + "[%d] is an array", index)
	// 				func_parse_array(intf.([]interface{}), level + 1)
	// 			}
	// 		}
	// 		return
	// 	}	// func_parse_array ends

	// 	func_parse_obj = func(obj map[string]interface{}, level int) {
	// 		prefix := strings.Repeat("    ", level)
	// 		for key, intf := range obj {
	// 			switch intf.(type) {
	// 			case string:
	// 				log.Infof(prefix + "\"%s\" = %s", key, intf.(string))
	// 			case float64:
	// 				integer := int64(intf.(float64))
	// 				if (float64(integer) == intf.(float64)) {
	// 					log.Infof(prefix + "\"%s\" = %d", key, integer)
	// 				} else {
	// 					log.Infof(prefix + "\"%s\" = %f", key, intf.(float64))
	// 				}
	// 			case bool:
	// 				log.Infof(prefix + "\"%s\" = %t", key, intf.(bool))
	// 			case map[string]interface{}:
	// 				log.Infof(prefix + "\"%s\" is an object", key)
	// 				func_parse_obj(intf.(map[string]interface{}), level + 1)
	// 			case []interface{}:
	// 				log.Infof(prefix + "\"%s\" is an array", key)
	// 				func_parse_array(intf.([]interface{}), level + 1)
	// 			}
	// 		}
	// 		return
	// 	}	// func_parse_obj ends

	// 	json_bytes := []byte("{\"str\": \"hello, json\", \"int\": 123, \"float\": 10.2, \"array\": [1, \"2\"], \"obj\": {\"num\": 10}, \"bool\": true, \"int64\": 4418489049307132905}")
	// 	var data map[string]interface{}
	// 	err := json.Unmarshal(json_bytes, &data)
	// 	if err != nil {
	// 		log.Error("unmarshal json error: " + err.Error())
	// 	} else {
	// 		func_parse_obj(data, 0)
	// 	}
	// }

	// start
	tars.Run()
}
