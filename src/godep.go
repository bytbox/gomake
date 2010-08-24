// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	. "container/vector"
	"fmt"
	"go/ast"
	"go/parser"
	"opts"
	"os"
	"path"
	"strings"
)

var showVersion = opts.LongFlag("version", "display version information")
var showNeeded = opts.Flag("n", "need", "display external dependencies")
var srcRoot = opts.Half("r", "root", "root directory of the source", "", "src")
var progName = "godep"

var roots = map[string]string{}

// prefix the root
func mkRoot(str string) string {
	return path.Join(*srcRoot, str)
}

func main() {
	opts.Usage = "[file1.go [...]]"
	opts.Description =
		`construct and print a dependency tree for the given source files.`
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
	// for each file, list dependencies
	for _, fname := range files {
		file, err := parser.ParseFile(fname, nil, parser.ImportsOnly)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		HandleFile(fname, file)
	}
	PrintAutoNotice()
	FindMain()
	if *showNeeded {
		PrintNeeded(".EXTERNAL: ", ".a")
	}
	// in any case, print as a comment
	PrintNeeded("# external packages: ", "")
	PrintDeps()
}

type Package struct {
	name     string
	files    *StringVector     // the files in this package
	packages map[string]string // dependencies
	hasMain  bool              // is this a main package with a `main` function
	path     string            // the path to this package
}

// packages is a mapping of package names (strings) to Package objects
var packages = map[string]Package{}

// FindMain finds all files which are in package 'main' and have a 'main'
// function.
func FindMain() {
	// for each file in the main package
	if pkg, ok := packages["main"]; ok {
		for _, fname := range *pkg.files {
			file, _ := parser.ParseFile(fname, nil, 0)
			ast.Walk(&MainCheckVisitor{fname}, file)
		}
	}
}

// PrintNeeded prints out a list of external dependencies to standard output.
func PrintNeeded(pre, ppost string) {
	// dependencies already displayed
	done := map[string]bool{}
	// start the list
	fmt.Print(pre)
	// for each package
	for _, pkg := range packages {
		// print all packages for which we don't have the source
		for _, pkgname := range pkg.packages {
			if _, ok := packages[pkgname]; !ok && !done[pkgname] {
				fmt.Printf("%s%s ", pkgname, ppost)
				done[pkgname] = true
			}
		}
	}
	fmt.Print("\n")
}

// PrintDeps prints out the dependency lists to standard output.
func PrintDeps() {
	// for each package
	for pkgname, pkg := range packages {
		if pkgname != "main" {
			// start the list
			fmt.Printf("%s.a: ", mkRoot(pkgname))
			// print all the files
			for _, fname := range *pkg.files {
				fmt.Printf("%s ", fname)
			}
			// print all packages for which we have the source
			// exception: if -n was supplied, print all packages
			for _, pkgname := range pkg.packages {
				_, ok := packages[pkgname]
				if ok || *showNeeded {
					fmt.Printf("%s.a ", mkRoot(pkgname))
				}
			}
			fmt.Printf("\n")
		}
	}
	common := StringVector{}
	// for the main package
	if main, ok := packages["main"]; ok {
		// consider all files not found in 'roots' to be common to
		// everything in this package
		for _, fname := range *main.files {
			if app, ok := roots[fname]; ok {
				fmt.Printf("%s: %s.${O}\n", app, app)
			} else {
				common.Push(fname)
			}
		}
		// for every application root
		for _, fname := range *main.files {
			if app, ok := roots[fname]; ok {
				// dependencies already displayed
				done := map[string]bool{}
				// print the file
				fmt.Printf("%s.${O}: %s ", app, fname)
				// print the common files
				for _, cfile := range common {
					fmt.Printf("%s ", cfile)
				}
				// print all packages for which we have the
				// source, or, if -n was supplied, print all
				for _, pkgname := range main.packages {
					_, ok := packages[pkgname]
					if ok || (*showNeeded && !done[pkgname]) {
						fmt.Printf("%s.a ", mkRoot(pkgname))
						done[pkgname] = true
					}
				}
				fmt.Printf("\n")
			}
		}
	}
}

func HandleFile(fname string, file *ast.File) {
	pkgname := file.Name.Name
	if pkg, ok := packages[pkgname]; ok {
		pkg.files.Push(fname)
	} else {
		packages[pkgname] = Package{
			files:    &StringVector{},
			packages: map[string]string{},
			hasMain:  false,
		}
		packages[pkgname].files.Push(fname)
	}
	ast.Walk(&ImportVisitor{packages[pkgname]}, file)
}

//
// ImportVisitor
//
// Finds a lists all imports for the scanned file.
//

type ImportVisitor struct {
	pkg Package
}

func (v ImportVisitor) Visit(node interface{}) ast.Visitor {
	// check the type of the node
	if spec, ok := node.(*ast.ImportSpec); ok {
		ppath := path.Clean(strings.Trim(string(spec.Path.Value), "\""))
		if _, ok = v.pkg.packages[ppath]; !ok {
			v.pkg.packages[ppath] = ppath
		}
	}
	return v
}

//
// MainCheckVisitor
//
// Used to check for a function named 'main' (usually in a package named 
// 'main'), to decide if a package should be made into its own executable.
//

type MainCheckVisitor struct {
	fname string
}

func addRoot(filename string) {
	fparts := strings.Split(filename, ".", -1)
	basename := fparts[0]
	roots[filename] = basename
}

func (v MainCheckVisitor) Visit(node interface{}) ast.Visitor {
	if decl, ok := node.(*ast.FuncDecl); ok {
		if decl.Name.Name == "main" {
			addRoot(v.fname)
		}
	}
	return v
}
