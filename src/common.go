package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var version = "0.2.0"

// Show version information
func ShowVersion() {
	fmt.Printf("%s (GoMake) v%s\n", progName, version)
}

func ReadConfig(filename string) (config map[string]string, err os.Error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	lines := strings.Split(string(content), "\n", -1)
	for line := range lines {
		parts := strings.Split(string(line), "=", 2)
		key, value := strings.Trim(parts[0], " "), 
			strings.Trim(parts[1], " ")
		config[key] = value
	}
	return
}
