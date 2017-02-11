package main

import (
	"flag"
	"go/ast"
	"go/token"
	"go/parser"
	"log"
	. "github.com/nirasan/ast-walker/lib"
	"regexp"
	"fmt"
	"reflect"
	"os"
)

var (
	filename = flag.String("f", "", "target file name.")
	pattern = flag.String("r", "", "regex pattern for search command.")
	typename = flag.String("t", "", "type name for search node.")
	found = false
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("ast-walker: ")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(2)
	}

	fset := token.NewFileSet()
	f, e := parser.ParseFile(fset, *filename, nil, 0)
	if e != nil {
		log.Fatal(e)
	}

	if *pattern == "" && *typename == "" {
		ast.Print(fset, f)
		return
	}

	if *pattern != "" {
		r, e := regexp.Compile(*pattern)
		if e != nil {
			log.Fatal(e)
		}
		Inspect(f, func(n ast.Node, h *History) bool {
			if h.MatchRegex(r) {
				print(fset, n, h)
			}
			return true
		})
	} else if *typename != "" {
		Inspect(f, func(n ast.Node, h *History) bool {
			t := reflect.TypeOf(n)
			if t != nil && t.String() == *typename {
				print(fset, n, h)
			}
			return true
		})
	}

	if !found {
		fmt.Println("Not found.")
	}
}

func print(fset *token.FileSet, n ast.Node, h *History) {
	if !found {
		found = true
		fmt.Println("Found!")
		fmt.Println("")
	}
	fmt.Println("ACCESSER_COMMAND: ")
	fmt.Println("    " + h.Path())
	fmt.Println("AST_PRINT: ")
	ast.Print(fset, n)
	fmt.Println("")
}
