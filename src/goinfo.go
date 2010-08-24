// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"opts"
	"os"
	"path"
	"strings"
)

var progName = "goinfo"

var showVersion = opts.LongFlag("version", "display version information")
var srcRoot = opts.Half("r", "root", "root directory of the source", "", "src")

// prefix the root
func mkRoot(str string) string {
	return path.Join(*srcRoot, str)
}

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
	GetPackageList()
	PrintAutoNotice()
	PrintFList()
	PrintPList()
	fmt.Println("GOPACKAGES = ${GOPKGS:=.${O}}")
	fmt.Println("GOARCHIVES = ${GOPKGS:=.a}")
}

// Print list of files
func PrintFList() {
	fmt.Printf("GOFILES = ")
	for _, fname := range files {
		fmt.Print(fname + " ")
	}
	fmt.Print("\n")
}

var packages = map[string]*struct{}{}

func GetPackageList() {
	for _, fname := range files {
		file, err := parser.ParseFile(fname, nil, parser.PackageClauseOnly)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		pname := file.Name.Name
		if pname == "main" {
			fullfile, err := parser.ParseFile(fname, nil, 0)
			if err == nil {
				v := &MainCheckVisitor{fname: fname}
				ast.Walk(v, fullfile)
				if v.hasMain {
					// get the name from the filename
					fparts := strings.Split(fname, ".", -1)
					basename := path.Base(fparts[0])
					packages[basename] = nil
				} else {
					packages[pname] = nil
				}
			}
		} else {
			packages[file.Name.Name] = nil
		}
	}
}

type MainCheckVisitor struct {
	fname   string
	hasMain bool
}

func (v *MainCheckVisitor) Visit(node interface{}) ast.Visitor {
	if decl, ok := node.(*ast.FuncDecl); ok {
		if decl.Name.Name == "main" {
			v.hasMain = true
		}
	}
	return v
}

// Print list of packages
func PrintPList() {
	fmt.Print("GOPKGS = ")
	for pname, _ := range packages {
		fmt.Printf("%s ", mkRoot(pname))
	}
	fmt.Print("\n")
}
