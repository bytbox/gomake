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

var showVersion = opts.Longflag("version", "display version information")
var outputFilename = opts.Shortopt("o", "file to write makefile to", "Makefile")

func main() {
	// if any arguments were given, this is being used as 'make'
	asMake := false
	for num, arg := range os.Args {
		if arg[0] != '-' && num != 0 {
			asMake = true
		}
	}
	if asMake {
		make, _ := exec.LookPath("make")
		os.Exec(make, os.Args, nil)
	}
	// parse and handle options
	opts.Parse()
	if *showVersion {
		ShowVersion()
		os.Exit(0)
	}
}
