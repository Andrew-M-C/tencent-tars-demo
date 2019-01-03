package main

import (
	"fmt"

	"github.com/TarsCloud/TarsGo/tars"

	"amc"
)

func main() {
	comm := tars.NewCommunicator()
	obj := fmt.Sprintf("amc.GoTarsServer.GoObj@tcp -h 127.0.0.1 -p 10015 -t 60000")
	app := new(amc.Go)
	comm.StringToProxy(obj, app)
	var out, i int32
	i = 123
	ret, err := app.Add(i, i*2, &out)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret, out)
}
