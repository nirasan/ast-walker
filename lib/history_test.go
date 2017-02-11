package ast_walker

import (
	"testing"
	"go/ast"
	"go/parser"
	"go/token"
	"bytes"
)

func TestHistory_Match(t *testing.T) {
	src := `
package main
type struct1 struct {
	field1 int
	field2 string
	field3 bool "json:\"filed3\""
	field4 struct {
		inner1 float32
		inner2 uint
	}
}
`

	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "", src, 0)

	//ast.Print(fset, f)

	Inspect(f, func(n ast.Node, history *History) bool {
		b1 := new(bytes.Buffer)
		switch x := n.(type) {
		case *ast.Field:
			ast.Fprint(b1, fset, x, ast.NotNilFilter)
		}

		b2 := new(bytes.Buffer)
		// struct fields path example:
		//   Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[0]
		if history.MatchString(`Fields.List\[\d+\]$`) {
			ast.Fprint(b2, fset, n, ast.NotNilFilter)
		}

		if bytes.Compare(b1.Bytes(), b2.Bytes()) != 0 {
			t.Errorf("not match:\ntype match:\n%s\nregex match:\n%s", b1.String(), b2.String())
		} else {
			t.Logf("match:\ntype match:\n%s\nregex match:\n%s", b1.String(), b2.String())

		}

		return true
	})
}