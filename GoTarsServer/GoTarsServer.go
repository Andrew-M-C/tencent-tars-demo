package main

import (
	"github.com/TarsCloud/TarsGo/tars"
	amc "amc/GoTarsServer/Amc"
)

func main() {
	imp := new(GoImp)
	app := new(amc.DateTime)
	cfg := tars.GetServerConfig()
	app.AddServant(imp, cfg.App + "." + cfg.Server + ".GoTarsObj") //Register Servant
	tars.Run()
}
