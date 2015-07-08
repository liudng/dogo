// Copyright 2014 The dogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/zhgo/console"
	"strings"
)

func main() {
	var c string
	flag.StringVar(&c, "c", console.WorkingDir+"/dogo.json", "Usage: dogo -c=/path/to/dogo.json")
	flag.Parse()

	var dogo Dogo

	gopath := console.Getenv("GOPATH")
	c = strings.Replace(c, "{GOPATH}", gopath, -1)
	r := map[string]string{"{GOPATH}": gopath}
    err := console.NewConfig(c).Replace(r).Parse(&dogo)
    if err != nil {
		fmt.Printf("[dogo] Warning: no configuration file loaded.\n")
	} else {
		fmt.Printf("[dogo] Loaded configuration file:\n")
		fmt.Printf("       %s\n", c)
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
