// Copyright 2014 The dogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/zhgo/kernel"
	"github.com/zhgo/config"
	"github.com/zhgo/console"
	"strings"
)

func main() {
	var c string
	flag.StringVar(&c, "c", kernel.WorkingDir+"/dogo.json", "Usage: dogo -c=/path/to/dogo.json")
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
		fmt.Println("[dogo] Warning: no configuration file loaded.")
	} else {
		fmt.Println("[dogo] Loaded configuration file: ", c)
	}

	dogo.NewMonitor()

	l := len(dogo.Files)
	if l > 0 {
		fmt.Printf("[dogo] Ready. %d files to be monitored.\n\n", l)

		dogo.BuildAndRun()

		dogo.Monitor()
	} else {
		fmt.Printf("[dogo] Error: Did not find any files. Press any key to exit.\n\n")
		var a string
		fmt.Scanf("%s", &a)
	}
}
