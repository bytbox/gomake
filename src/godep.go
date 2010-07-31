package main

import (
	. "container/vector"
	"fmt"
	"go/ast"
	"go/parser"
	"opts"
	"os"
	"strings"
)

var showVersion = opts.Longflag("version", "display version information")
var configFilename = opts.Longopt("config", 
	"name of configuration file", "godep.cfg")
var progName = "godep"
var config map[string]string

func main() {
	opts.Usage("file1.go [...]")
	opts.Description(`construct and print a dependency tree for the given source files.`)
	// parse and handle options
	opts.Parse()
	if *showVersion {
		ShowVersion()
		os.Exit(0)
	}
	// read the configuration
	readConfig()
	// for each file, list dependencies
	for _, fname := range opts.Args {
		file, err := parser.ParseFile(fname, nil, nil, parser.ImportsOnly)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		HandleFile(fname, file)
	}
	PrintDeps()
}

var roots = map[string]string{}

func readConfig() {
	// read from the configuration file, if any (discard the error)
	config, _ = ReadConfig(*configFilename)
	rootstr, ok := config["roots"]
	if !ok {
		return
	}
	// extract the roots
	rootlist := strings.Split(rootstr, " ", -1)
	for _, root := range rootlist {
		kv := strings.Split(strings.Trim(root, " "),
			":", 2)
		if len(kv) == 2 {
			roots[kv[1]] = kv[0]
		}
	}
}

type Package struct {
	// The files on which this package depends.
	files *StringVector
	// The packages on which this package depends.
	packages *StringVector
}

var packages = map[string]Package{}

// PrintDeps prints out the dependency lists to standard output.
func PrintDeps() {
	// for each package
	for pkgname, pkg := range packages {
		if pkgname != "main" {
			// start the list
			fmt.Printf("%s.${O}: ", pkgname)
			// print all the files
			for _, fname := range *pkg.files {
				fmt.Printf("%s ", fname)
			}
			// print all packages for which we have the source
			for _, pkgname := range *pkg.packages {
				if _, ok := packages[pkgname]; ok {
					fmt.Printf("%s.${O} ", pkgname)
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
		for _, fname := range *main.files {
			if app, ok := roots[fname]; ok {
				fmt.Printf("%s.${O}: %s ", app, fname)
				for _, cfile := range common {
					fmt.Printf("%s ",cfile)
				}
				fmt.Printf("\n")
			}
		}
	}
}

func HandleFile(fname string, file *ast.File) {
	pkgname := file.Name.Name()
	if pkg, ok := packages[pkgname]; ok {
		pkg.files.Push(fname)
	} else {
		packages[pkgname] = Package{&StringVector{}, &StringVector{}}
		packages[pkgname].files.Push(fname)
	}
	ast.Walk(&Visitor{packages[pkgname]}, file)
}

type Visitor struct {
	pkg Package
}

func (v Visitor) Visit(node interface{}) ast.Visitor {
	// check the type of the node
	if spec, ok := node.(*ast.ImportSpec); ok {
		path := strings.Trim(string(spec.Path.Value), "\"")
		v.pkg.packages.Push(path)
	}
	return v
}
