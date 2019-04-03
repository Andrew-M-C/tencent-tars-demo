package dynamic
import (
	"github.com/Andrew-M-C/go-tools/log"
	"github.com/Andrew-M-C/go-tools/jsonconv"
	"github.com/Andrew-M-C/go-tools/str"
	"github.com/satori/go.uuid"
	qrcode "github.com/skip2/go-qrcode"
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
		Name: "ticket",
		Value: "",
		Expires: time.Now(),
		Domain: "andrewmc.cn",
		Path: "/",
	})
	return
}


func redirectToLoginPage(w http.ResponseWriter, r *http.Request) {
	log.Debug("redirect to login page")
	http.Redirect(w, r, "/html/login.html", 302)
	return
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	info := httpparser.GetHttpRequestInfo(r)
	url := info.Url
	log.Debug("request URL: %s", url)
	log.Debug("client IP: %s", info.IP)
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

	// generate a ticket
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
		Name: "ticket",
		Value: t.Session,
		Expires: t.Expired,
		Domain: "andrewmc.cn",
		Path: "/",
	})

	// check redirect_url
	resp := jsonconv.NewObject()
	resp.Set(jsonconv.NewInt(0), "code")

	resp_str, _ := resp.Marshal()
	w.Write([]byte(resp_str))
	return
}


func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	// check cookies
	uid_in_cookie, _ := r.Cookie("uid")
	ticket_in_cookie, _ := r.Cookie("ticket")
	if uid_in_cookie == nil || ticket_in_cookie == nil {
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
	if ticket.Session != ticket_in_cookie.Value {
		log.Info("ticket '%s' mismatch", ticket_in_cookie.Value)
		clearCookieSession(w)
		redirectToLoginPage(w, r)
		return
	}
	if time.Now().After(ticket.Expired) {
		log.Info("ticket '%s' expired", ticket_in_cookie.Value)
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

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Logout")
	uid_in_cookie, _ := r.Cookie("uid")
	ticket_in_cookie, _ := r.Cookie("ticket")
	if uid_in_cookie == nil || ticket_in_cookie == nil {
		log.Info("this is an invalid request for logout, simply redirect it")
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
	if ticket.Session == ticket_in_cookie.Value {
		log.Info("ticket '%s' match, clear it", ticket_in_cookie.Value)
		delete(_tickets, uid_in_cookie.Value)
		clearCookieSession(w)
		redirectToLoginPage(w, r)
		return
	}
	clearCookieSession(w)
	redirectToLoginPage(w, r)
	return
}


func TicketValidateHandler(w http.ResponseWriter, r *http.Request) {
	info := httpparser.GetHttpRequestInfo(r)
	fail_resp := []byte(`{"code":-1,"data":{"result":false}}`)
	succ_resp := []byte(`{"code":0,"data":{"result":true}}`)
	log.Debug("client IP: %s", info.IP)

	for key, val := range info.Query {
		log.Debug("Query - [%s] = '%s'", key, val)
	}

	// body, _ := ioutil.ReadAll(r.Body)
	// if len(body) > 0 {
	// 	log.Debug("Get request: %s", string(body))
	// }

	var uid_param, ticket_param string
	uid_in_cookie, _ := r.Cookie("uid")
	ticket_in_cookie, _ := r.Cookie("ticket")
	if uid_in_cookie == nil || ticket_in_cookie == nil {
		// let's try query parameter
		uid_in_query, exist_uid := info.Query["uid"]
		ticket_in_query, exist_ticket := info.Query["ticket"]
		if exist_uid && exist_ticket {
			uid_param = uid_in_query
			ticket_param = ticket_in_query
		} else {
			log.Info("this is an invalid request for logout, simply redirect it")
			clearCookieSession(w)
			w.Write(fail_resp)
			return
		}
	} else {
		uid_param = uid_in_cookie.Value
		ticket_param = ticket_in_cookie.Value
	}

	_lock.Lock()
	defer _lock.Unlock()
	ticket, exist := _tickets[uid_param]
	if false == exist {
		log.Info("ticket for '%s' not exist", uid_param)
		clearCookieSession(w)
		w.Write(fail_resp)
		return
	}
	if ticket.Session != ticket_param {
		log.Info("ticket '%s' not match, clear it", ticket_param)
		clearCookieSession(w)
		w.Write(fail_resp)
		return
	}

	log.Debug("This is a valid cookie")
	w.Write(succ_resp)
	return
}


func QRImageHandler(w http.ResponseWriter, r *http.Request) {
	info := httpparser.GetHttpRequestInfo(r)
	req_url := info.Url
	log.Debug("request URL: %s", req_url)
	defer w.Header().Set("Content-Type", "image/png")

	// test a QR image generator
	qr, err := qrcode.Encode("https://www.google.com", qrcode.Medium, 256)
	if err != nil {
		log.Error("Failed to generate QR image: %s", err.Error())
		return
	}
	w.Write(qr)
	return
}
