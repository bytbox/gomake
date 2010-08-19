// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var version = "0.2.2"

// Show version information
func ShowVersion() {
	fmt.Printf("%s (GoMake) v%s\n", progName, version)
}

func ReadConfig(filename string) (map[string]string, os.Error) {
	config := map[string]string{}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
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
