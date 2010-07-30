package main

import "fmt"
import (
	"go/ast"
	"go/parser"
	"opts"
	"os"
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
		file, _ := parser.ParseFile(fname, nil, nil, parser.ImportsOnly)
		fmt.Printf("%s is in package %s\n", fname, file.Name.Name())
		ast.Walk(&Visitor{},file)
	}
}

// Show version information
func ShowVersion() {
	fmt.Printf("godep v%s\n",version)
}

type Visitor struct {}

func (v Visitor) Visit(node interface{}) ast.Visitor {
	// check the type of the node
	if spec, ok := node.(*ast.ImportSpec); ok {
		fmt.Printf("importing %s\n",spec.Path.Value)
	}
	return v
}

