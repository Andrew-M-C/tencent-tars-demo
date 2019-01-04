package main

import (
	"github.com/TarsCloud/TarsGo/tars"
	Amc "amc/GoTarsServer/Amc"
)

func main() {
	imp := new(GoImp)
	app := new(Amc.Go)
	cfg := tars.GetServerConfig()
	app.AddServant(imp, cfg.App + "." + cfg.Server + ".GoTarsObj") //Register Servant
	tars.Run()
}
