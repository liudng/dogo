// Copyright 2014 The dogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/favframework/config"
	"github.com/favframework/console"
	"strings"
)

var WorkingDir string = config.WorkingDir()

func main() {
	var c string
	flag.StringVar(&c, "c", WorkingDir+"/dogo.json", "Usage: dogo -c=/path/to/dogo.json")
	flag.Parse()

	New(c)
}

func New(c string) {
	var dogo Dogo

	gopath := console.Getenv("GOPATH")

	r := make(map[string]string)
	r["{GOPATH}"] = gopath

	c = strings.Replace(c, "{GOPATH}", gopath, -1)

	err := config.LoadJSONFile(&dogo, c, r)
	if err != nil {
		fmt.Println("[dogo] No config file loaded.")
	} else {
		fmt.Println("[dogo] Load config file: ", c)
	}

	dogo.NewMonitor()
	dogo.BuildAndRun()
	dogo.Monitor()
}
