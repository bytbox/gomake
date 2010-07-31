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

func ReadConfig(filename string) (map[string]string, os.Error) {
	config := map[string]string{}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n", -1)
	for _, line := range lines {
		parts := strings.Split(string(line), "=", 2)
		if len(parts) == 2 {
			key, value := strings.Trim(parts[0], " "), 
				strings.Trim(parts[1], " ")
			config[key] = value
		}
	}
	return config, err
}
