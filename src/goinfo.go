// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"opts"
	"os"
	"path"
)

var progName = "goinfo"

var showVersion = opts.LongFlag("version", "display version information")

func main() {
	// parse and handle options
	opts.Parse()
	if *showVersion {
		ShowVersion()
		os.Exit(0)
	}
	// if there are no files, generate a list
	if len(opts.Args) == 0 {
		path.Walk(".", GoFileFinder{}, nil)
	} else {
		for _, fname := range opts.Args {
			files.Push(fname)
		}
	}
	PrintFList()
	PrintPList()
}

// Print list of files
func PrintFList() {
	fmt.Printf("GOFILES = ")
	for _, fname := range files {
		fmt.Print(fname+" ")
	}
	fmt.Print("\n")
}

// Print list of packages
func PrintPList() {
}
