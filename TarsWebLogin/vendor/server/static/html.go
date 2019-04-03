package static
import (
	"github.com/Andrew-M-C/go-tools/log"
	"utils/httpparser"
	"utils/config"
	"net/http"
	"io/ioutil"
	"strings"
)

var _fileNotFoundPage = []byte(`
<head><title>404 Not Found</title></head>
<body bgcolor="white">
<center><h1>404 Not Found</h1></center>
</body>
</html>
`)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	info := httpparser.GetHttpRequestInfo(r)
	url := info.Url
	if url == "/" {
		http.Redirect(w, r, "html/login.html", 302)
	} else {
		// https://stackoverflow.com/questions/40096750/how-to-set-http-status-code-on-http-responsewriter
		w.WriteHeader(http.StatusNotFound)
	}
	return
}

func HtmlHandler(w http.ResponseWriter, r *http.Request) {
	info := httpparser.GetHttpRequestInfo(r)
	url := info.Url
	log.Debug("request URL: %s", url)
	log.Debug("client IP: %s", info.IP)

	if url != strings.ToLower(url) {
		http.Redirect(w, r, strings.ToLower(url), 302)
	}

	if strings.HasSuffix(url, ".html") {
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
	} else if strings.HasSuffix(url, ".ico") {
		w.Header().Set("Content-Type", "image/x-icon")
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	// open file and send back
	file, err := ioutil.ReadFile(config.GetHomeDir() + url)
	if err != nil {
		log.Error("Cannot read file '%s': %s", url, err.Error())
		w.Write(_fileNotFoundPage)
	} else {
		w.Write(file)
	}
	return
}


func JsHandler(w http.ResponseWriter, r *http.Request) {
	info := httpparser.GetHttpRequestInfo(r)
	url := info.Url
	log.Debug("request URL: %s", url)

	if url != strings.ToLower(url) {
		http.Redirect(w, r, strings.ToLower(url), 302)
	}

	w.Header().Set("Content-Type", "text/javascript;charset=utf-8")
	// open file and send back
	file, err := ioutil.ReadFile(config.GetHomeDir() + url)
	if err != nil {
		log.Error("Cannot read file '%s': %s", url, err.Error())
		w.Write(_fileNotFoundPage)
	} else {
		w.Write(file)
	}
	return
}
