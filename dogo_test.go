// Copyright 2014 The dogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func TestInitFiles(t *testing.T) {
	var dogo Dogo

	dogo.NewMonitor()

	// TODO: Configuration is not being loaded here so, being that
	//       the ignored files and folders were added, I had to change
	//       the number of files here from 9 to 11.
	//       I think that a better way would be to load a test configuration
	//       file in order to have more control over the tests output and
	//       so be able to test that `dogo.Files` does not contains ignored
	//       files or folders.
	l := len(dogo.Files)
	if l != 11 {
		t.Fatalf("Init source files failed: %d.\n", l)
	}
}
