package main

import (
	"github.com/favframework/config"
)

var WorkingDir string = config.WorkingDir()

func main(){
	//log.Printf("%s\n", WorkingDir)

	var dogo Dogo

	config.LoadJSONFile(&dogo, WorkingDir+"/dogo.json", nil)

	dogo.NewMonitor()
}
