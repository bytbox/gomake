package main

import (
	"fmt"
)

var version = "0.2.0"

// Show version information
func ShowVersion() {
	fmt.Printf("%s (GoMake) v%s\n", progName, version)
}
