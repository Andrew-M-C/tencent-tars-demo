package httpparser
import (
	// "github.com/Andrew-M-C/go-tools/log"
	"net/http"
	// "fmt"
	"strconv"
	"strings"
)

type HttpRequestInfo struct {
	Ip			string
	Port		int
	Host		string
	Url			string
	Header		map[string]string
	Query		map[string]string
}

func GetHttpRequestInfo(r *http.Request) *HttpRequestInfo {
	info := new(HttpRequestInfo)
	info.Header = make(map[string]string)
	info.Query  = make(map[string]string)

	// IP and port
	info.Ip = r.Header.Get("X-Real-IP")
	port_str := r.Header.Get("X-Real-Port")
	port, err := strconv.Atoi(port_str)
	if nil != err {
		port = -1
	}
	info.Port = port

	// Host
	info.Host = r.Host

	// URL
	full_url := r.RequestURI
	question_mark_index := strings.IndexAny(full_url, "?")
	if question_mark_index >= 0 {
		info.Url = full_url[0:question_mark_index]
	} else {
		info.Url = full_url
	}

	// header
	for key, value := range r.Header {
		// log.Debug(fmt.Sprintf("header [%s] - %s", key, value[0]))
		info.Header[key] = value[0]
	}

	// query param
	r.ParseForm()
	for key, value := range r.Form {
		// log.Debug(fmt.Sprintf("param [%s] - %s", key, value[0]))
		info.Query[key] = value[0]
	}

	return info
}
