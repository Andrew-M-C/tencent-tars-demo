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
	http.Redirect(w, r, "html/login.html", 302)
	return
}

func HtmlHandler(w http.ResponseWriter, r *http.Request) {
	info := httpparser.GetHttpRequestInfo(r)
	url := info.Url
	log.Debug("request URL: %s", url)

	if url != strings.ToLower(url) {
		http.Redirect(w, r, strings.ToLower(url), 302)
	}

	w.Header().Set("Content-Type", "text/html;charset=utf-8")
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
