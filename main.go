package main

import (
	"flag"
	"fmt"
	. "github.com/nirasan/ast-walker/lib"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"reflect"
	"regexp"
)

var (
	command = flag.String("command", "", "command for search.")
	pattern = flag.String("regex", "", "regex pattern for search command.")
	typename = flag.String("type", "", "type name for search node.")
	found    = false
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tShow ast.Print:\n")
	fmt.Fprintf(os.Stderr, "\t\tast-walker FILENAME\n")
	fmt.Fprintf(os.Stderr, "\tShow command result:\n")
	fmt.Fprintf(os.Stderr, "\t\tast-walker -command COMMAND FILENAME\n")
	fmt.Fprintf(os.Stderr, "\tShow command matched result:\n")
	fmt.Fprintf(os.Stderr, "\t\tast-walker -regex REGEX FILENAME\n")
	fmt.Fprintf(os.Stderr, "\tShow type mached result:\n")
	fmt.Fprintf(os.Stderr, "\t\tast-walker -type NAME FILENAME\n")
	fmt.Fprintf(os.Stderr, "Hint:\n")
	fmt.Fprintf(os.Stderr, "\tSearch first struct declaration by command:\n")
	fmt.Fprintf(os.Stderr, "\t\tast-walker -command %s FILENAME\n", `Decls[0].(*ast.GenDecl).Specs[0]`)
	fmt.Fprintf(os.Stderr, "\tSearch struct declaration by type:\n")
	fmt.Fprintf(os.Stderr, "\t\tast-walker -type %s FILENAME\n", `*ast.TypeSpec`)
	fmt.Fprintf(os.Stderr, "\tSearch struct field declaration by regex:\n")
	fmt.Fprintf(os.Stderr, "\t\tast-walker -regex %s FILENAME\n", `'Fields.List\[\d+\]$'`)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("ast-walker: ")
	flag.Usage = Usage
	flag.Parse()

	filename := flag.Arg(0)

	if filename == "" {
		flag.Usage()
		os.Exit(2)
	}

	fset := token.NewFileSet()
	f, e := parser.ParseFile(fset, filename, nil, 0)
	if e != nil {
		log.Fatal(e)
	}

	if *command == "" && *pattern == "" && *typename == "" {
		ast.Print(fset, f)
		return
	}

	if *command != "" {
		Inspect(f, func(n ast.Node, h *History) bool {
			if h.Path() == *command {
				print(fset, n, h)
			}
			return true
		})
	} else if *pattern != "" {
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
	fmt.Println("COMMAND: ")
	fmt.Println("    " + h.Path())
	fmt.Println("AST_PRINT: ")
	ast.Print(fset, n)
	fmt.Println("")
}
