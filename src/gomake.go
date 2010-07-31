package main

import (
	"opts"
	"os"
)

var progName = "gomake"

var showVersion = opts.Longflag("version", "display version information")
var outputFilename = opts.Shortopt("o", "file to write makefile to", "Makefile")

func main() {
	// parse and handle options
	opts.Parse()
	if *showVersion {
		ShowVersion()
		os.Exit(0)
	}
}
