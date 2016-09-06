// Copyright 2014 The dogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (d *Dogo) Monitor() {
	for {
		d.Compare()

		if d.isModified == true {
			d.BuildAndRun()
		}

		time.Sleep(time.Duration(1 * time.Second))
	}
}

//compare source file's modify time
func (d *Dogo) Compare() {
	changed := false

	for p, t := range d.Files {
		info, err := os.Stat(p)
		if err != nil {
			delete(d.Files, p)
			continue
		}

		//new modtime
		nt := info.ModTime()

		if nt.Sub(t) > 0 {
			d.Files[p] = nt
			changed = true
			fmt.Printf("[dogo] Changed files: %s\n", filepath.Base(p))
		}
	}

	if changed == true {
		d.isModified = true
	} else {
		d.isModified = false
	}
}
