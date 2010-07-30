package main

import "fmt"
import (
	. "container/vector"
	"go/ast"
	"go/parser"
	"opts"
	"os"
	"strings"
)

var version = "0.0.1"

var showVersion = opts.Longflag("version", "display version information")

func main() {
	opts.Usage("file1.go [...]")
	opts.Description(`construct and print a dependency tree for the given source files.`)
	// parse and handle options
	opts.Parse()
	if *showVersion {
		ShowVersion()
		os.Exit(0)
	}
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

// Show version information
func ShowVersion() {
	fmt.Printf("godep v%s\n", version)
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
