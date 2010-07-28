package main

import (
	"fmt"
	"go/parser"
	"opts"
	"os"
)

var version = "0.0.1"

var showVersion = opts.Longflag("version", "display version information")

func main() {
	// parse and handle options
	opts.Parse()
	if *showVersion {
		ShowVersion()
		os.Exit(0)
	}
	// for each file, list dependencies
	for _, fname := range opts.Args {
		fmt.Printf("%s\n",fname)
		file, _ := parser.ParseFile(fname, nil, nil, parser.ImportsOnly)
	}
}

// Show version information
func ShowVersion() {
	fmt.Printf("godep v%s\n",version)
}

