package main

import (
	"net/http"
	"github.com/TarsCloud/TarsGo/tars"
)

func dummyHttpRootHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func main() {
	{
		cfg := tars.GetServerConfig()
		mux := &tars.TarsHttpMux{}
		server := cfg.App + "." + cfg.Server
		mux.HandleFunc("/", dummyHttpRootHandler)
		tars.AddHttpServant(mux, server + ".ConsumerObj")
	}
	tars.Run()
}
