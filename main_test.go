package main

import (
	"testing"
	"go/ast"
	"go/parser"
	"go/token"
	"fmt"
	"strings"
)

func TestInspect(t *testing.T) {
	src := `
package main

const c1 = 1

const (
  c2 = iota
  c3
  c4
)

type mystruct struct {
  n1 uint
  n2 string
}
`

	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "", src, 0)

	ast.Print(fset, f)

	Inspect(f, func(n ast.Node, history []string) bool {
		switch x := n.(type) {
		case *ast.Ident:
			fmt.Printf("Path: %v\n", strings.Join(history, "."))
			ast.Print(fset, x)

		}
		return true
	})

	ast.Print(fset, f.Decls[2].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[1].Names[0])
}

func TestInspect2(t *testing.T) {
	src := `
package main

type mystruct struct {
  Arg1 uint "json:\"arg1\""
  Arg2 string
}
`

	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "", src, 0)

	ast.Print(fset, f)

	Inspect(f, func(n ast.Node, history []string) bool {
		switch x := n.(type) {
		case *ast.StructType:
			fmt.Printf("Path: %v\n", strings.Join(history, "."))
			ast.Print(fset, x)
		}
		return true
	})

	//ast.Print(fset, f.Decls[2].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[1].Names[0])
}