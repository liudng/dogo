package main

import (
	"github.com/favframework/config"
	"github.com/favframework/console"
	"flag"
)

var WorkingDir string = config.WorkingDir()

func main(){
	//log.Printf("%s\n", WorkingDir)

	var c string
	flag.StringVar(&c, "c", WorkingDir+"/dogo.json", "Usage: dogo -c=/path/to/dogo.json")
	flag.Parse()

	var dogo Dogo

	r := make(map[string]string)
	r["{GOPATH}"] = console.Getenv("GOPATH")

	config.LoadJSONFile(&dogo, c, r)

	dogo.NewMonitor()
}
