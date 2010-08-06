package main

import (
	"opts"
	"os"
)

var progName = "getgo"

var showVersion = opts.Longflag("version", "display version information")

func main() {
	// parse and handle options
	opts.Parse()
	if *showVersion {
		ShowVersion()
		os.Exit(0)
	}
}
