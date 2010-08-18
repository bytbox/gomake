// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"exec"
	"opts"
	"os"
)

var progName = "gomake"

var showVersion = opts.LongFlag("version", "display version information")
var outputFilename = opts.Single("o", "", "file to write makefile to", "Makefile")

func main() {
	// parse and handle options
	opts.Parse()
	if *showVersion {
		ShowVersion()
		os.Exit(0)
	}
}
