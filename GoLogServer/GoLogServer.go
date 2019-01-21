package main

import (
	"github.com/TarsCloud/TarsGo/tars"
	log "./Logf"
)

func main() {
	imp := new(GoLogObj)
	app := new(log.Log)
	cfg := tars.GetServerConfig()
	app.AddServant(imp, cfg.App + "." + cfg.Server + ".GoLogObj")
	tars.Run()
}
