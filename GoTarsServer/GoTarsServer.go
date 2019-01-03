package main

import (
	"github.com/TarsCloud/TarsGo/tars"

	"amc"
)

func main() { //Init servant
	imp := new(GoImp)                                    //New Imp
	app := new(amc.Go)                                 //New init the A Tars
	cfg := tars.GetServerConfig()                               //Get Config File Object
	app.AddServant(imp, cfg.App+"."+cfg.Server+".GoObj") //Register Servant
	tars.Run()
}
