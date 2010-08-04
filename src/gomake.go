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
	// parse and handle options
	opts.Parse()
	if *showVersion {
		ShowVersion()
		os.Exit(0)
	}
	// if any arguments were given, this is being used as 'make'
	if len(opts.Args) > 0 {
		make, _ := exec.LookPath("make")
		os.Exec(make,os.Args,nil)
	}
}
