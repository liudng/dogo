package main

import (
	"github.com/favframework/config"
	"github.com/favframework/console"
)

var WorkingDir string = config.WorkingDir()

func main(){
	//log.Printf("%s\n", WorkingDir)

	var dogo Dogo

	r := make(map[string]string)
	r["{GOPATH}"] = console.Getenv("GOPATH")

	config.LoadJSONFile(&dogo, WorkingDir+"/dogo.json", r)

	dogo.NewMonitor()
}
