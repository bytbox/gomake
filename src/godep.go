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

var showVersion = opts.Longflag("version", "display version information")
var progName = "godep"

var roots = map[string] string{}
var files = StringVector{}

type GoFileFinder struct {}
func (f GoFileFinder) VisitDir(path string, finfo *os.FileInfo) bool {
	return true
}

func (f GoFileFinder) VisitFile(fpath string, finfo *os.FileInfo) {
	if path.Ext(fpath) == ".go" {
		files.Push(fpath)
	}
}

func main() {
	opts.Usage("[file1.go [...]]")
	opts.Description(`construct and print a dependency tree for the given source files.`)
	// parse and handle options
	opts.Parse()
	if *showVersion {
		ShowVersion()
		os.Exit(0)
	}
	// if there are no files, generate a list
	if len(opts.Args) == 0 {
		path.Walk(".",GoFileFinder{}, nil)
	} else {
		for _, fname := range opts.Args {
			files.Push(fname)
		}
	}
	// for each file, list dependencies
	for _, fname := range files {
		file, err := parser.ParseFile(fname, nil, nil, parser.ImportsOnly)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		HandleFile(fname, file)
	}
	FindMain()
	PrintDeps()
}

type Package struct {
	files *StringVector
	packages *StringVector
	hasMain bool
}

var packages = map[string]Package{}

func FindMain() {
	// for each file in the main package
	if pkg, ok := packages["main"]; ok {
		for _, fname := range *pkg.files {
			file, _ := parser.ParseFile(fname, nil, nil, 0)
			ast.Walk(&MainCheckVisitor{fname},file)
		}
	}
}

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
		packages[pkgname] = Package{&StringVector{}, &StringVector{}, false}
		packages[pkgname].files.Push(fname)
	}
	ast.Walk(&ImportVisitor{packages[pkgname]}, file)
}

type ImportVisitor struct {
	pkg Package
}

func (v ImportVisitor) Visit(node interface{}) ast.Visitor {
	// check the type of the node
	if spec, ok := node.(*ast.ImportSpec); ok {
		path := strings.Trim(string(spec.Path.Value), "\"")
		v.pkg.packages.Push(path)
	}
	return v
}

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
		if decl.Name.Name() == "main" {
			addRoot(v.fname)
		}
	}
	return v
}
