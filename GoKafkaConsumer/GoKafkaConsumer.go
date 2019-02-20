package main

import (
	"net/http"
	"github.com/TarsCloud/TarsGo/tars"
	"../GoLogger/log"
)

func dummyHttpRootHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func main() {
	// configure tars
	cfg := tars.GetServerConfig()
	mux := &tars.TarsHttpMux{}
	server := cfg.App + "." + cfg.Server
	mux.HandleFunc("/", dummyHttpRootHandler)
	tars.AddHttpServant(mux, server + ".ConsumerObj")
	log.Infof("%s ready to run", server)

	// start actual object
	err := startConsumers(0)
	if err != nil {
		log.Error("Failed to start consumer: " + err.Error())
	}

	// start tars
	log.Infof("%s ready to run", server)
	tars.Run()
}
