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

	l := len(dogo.Files)
	if l != 7 {
		t.Fatalf("Init source files failed: %d.\n", l)
	}
}
