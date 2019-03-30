package dynamic
import (
	"github.com/Andrew-M-C/go-tools/log"
	"github.com/Andrew-M-C/go-tools/jsonconv"
	"github.com/Andrew-M-C/go-tools/str"
	"github.com/satori/go.uuid"
	"utils/httpparser"
	"utils/config"
	"net/http"
	"io/ioutil"
	"sync"
	"time"
	"fmt"
)

/* cookie:
uid_cookie:=&http.Cookie{
        Name:   "uid",
        Value:    uid,
        Path:     "/",
        HttpOnly: false,
        MaxAge:   maxAge
    }
 */

type ticket struct {
	Session		string
	Expired		time.Time
}

var (
	_lock	sync.Mutex
	_tickets = make(map[string]*ticket)
)


func clearCookieSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name: "uid",
		Value: "",
		Expires: time.Now(),
		Domain: "andrewmc.cn",
		Path: "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: "",
		Expires: time.Now(),
		Domain: "andrewmc.cn",
		Path: "/",
	})
	return
}


func redirectToLoginPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/html/login.html", 302)
	return
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	info := httpparser.GetHttpRequestInfo(r)
	url := info.Url
	log.Debug("request URL: %s", url)
	defer w.Header().Set("Content-Type", "application/json;charset=utf-8")

	// get request
	body_bytes, _ := ioutil.ReadAll(r.Body)
	body_str := string(body_bytes)
	log.Debug("Got request: %s", body_str)
	body, err := jsonconv.NewFromString(body_str)
	if err != nil {
		w.Write([]byte(`{"code":-1,"msg":"param error"}`))
		return
	}

	uid, _ := body.GetString("uid")
	hash, _ := body.GetString("hash")
	if str.Empty(uid) || str.Empty(hash) {
		w.Write([]byte(`{"code":-1,"msg":"param error"}`))
		return
	}
	log.Info(">> %s << login complete", uid)

	// get configurations from configuration (This is an example)
	log.Info("request login >> %s <<", uid)
	// conf_Str, _ := config.GetConfig().Marshal()
	// log.Info("config: %s", conf_Str)
	hash_in_conf, err := config.GetConfig().GetString("users", uid, "hash")
	if err != nil {
		log.Error("Failed to get %s: %s", uid, err.Error())
		w.Write([]byte(`{"code":-1,"msg":"user invalid"}`))
		return
	}
	if hash_in_conf != hash {
		log.Error("hash %s mismatch", hash)
		w.Write([]byte(`{"code":-1,"msg":"password mismatch"}`))
		return
	}

	// generate an session
	_lock.Lock()
	delete(_tickets, uid)
	t := new(ticket)
	u4, _ := uuid.NewV4()
	t.Session = u4.String()
	t.Expired = time.Now().AddDate(0, 0, 1)
	_tickets[uid] = t
	_lock.Unlock()
	log.Info("Update ticket for %s", uid)

	// login OK, get cookies
	http.SetCookie(w, &http.Cookie{
		Name: "uid",
		Value: uid,
		Expires: time.Now().AddDate(1, 0, 0),
		Domain: "andrewmc.cn",
		Path: "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: t.Session,
		Expires: t.Expired,
		Domain: "andrewmc.cn",
		Path: "/",
	})

	w.Write([]byte(`{"code":0}`))
	return
}


func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	// check cookies
	uid_in_cookie, _ := r.Cookie("uid")
	sess_in_cookie, _ := r.Cookie("session")
	if uid_in_cookie == nil || sess_in_cookie == nil {
		log.Info("this is an invalid request for success.html, redirect it")
		clearCookieSession(w)
		redirectToLoginPage(w, r)
		return
	}

	_lock.Lock()
	defer _lock.Unlock()
	ticket, exist := _tickets[uid_in_cookie.Value]
	if false == exist {
		log.Info("ticket for '%s' not exist", uid_in_cookie.Value)
		clearCookieSession(w)
		redirectToLoginPage(w, r)
		return
	}
	if ticket.Session != sess_in_cookie.Value {
		log.Info("ticket '%s' mismatch", sess_in_cookie.Value)
		clearCookieSession(w)
		redirectToLoginPage(w, r)
		return
	}
	if sess_in_cookie.Expires.After(ticket.Expired) {
		log.Info("ticket '%s' expired", sess_in_cookie.Value)
		clearCookieSession(w)
		redirectToLoginPage(w, r)
		return
	}

	log.Info("Session for '%s' valid", uid_in_cookie.Value)

	// read success.html
	file, err := ioutil.ReadFile(config.GetHomeDir() + "/html/success.html")
	if err != nil {
		log.Error("Cannot read success.html")
		clearCookieSession(w)
		redirectToLoginPage(w, r)
		return
	} else {
		full_html := fmt.Sprintf(string(file), uid_in_cookie.Value)
		w.Write([]byte(full_html))
		return
	}
}
