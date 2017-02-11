package main

import (
	"fmt"
	"go/ast"
	"reflect"
)

type Visitor interface {
	Visit(node ast.Node, history []string) (w Visitor)
}

func walkIdentList(v Visitor, list []*ast.Ident, history []string, name string) {
	for i, x := range list {
		Walk(v, x, AddListHistory(history, name, i), false)
	}
}

func walkExprList(v Visitor, list []ast.Expr, history []string, name string) {
	for i, x := range list {
		Walk(v, x, AddListHistory(history, name, i), true)
	}
}

func walkStmtList(v Visitor, list []ast.Stmt, history []string, name string) {
	for i, x := range list {
		Walk(v, x, AddListHistory(history, name, i), true)
	}
}

func walkDeclList(v Visitor, list []ast.Decl, history []string, name string) {
	for i, x := range list {
		Walk(v, x, AddListHistory(history, name, i), true)
	}
}

func AddListHistory(src []string, name string, i int) []string {
	return AddHistory(src, fmt.Sprintf("%s[%d]", name, i))
}

func AddHistory(src []string, name string) []string {
	dst := make([]string, len(src)+1)
	copy(dst, src)
	dst[len(dst)-1] = name
	return dst
}

func Walk(v Visitor, node ast.Node, history []string, parentIsInterface bool) {
	if v = v.Visit(node, history); v == nil {
		return
	}

	if parentIsInterface {
		history = AddHistory(history, fmt.Sprintf("(%s)", reflect.TypeOf(node)))
	}

	switch n := node.(type) {
	case *ast.Comment:

	case *ast.CommentGroup:
		for i, c := range n.List {
			Walk(v, c, AddListHistory(history, "List", i), false)
		}

	case *ast.Field:
		if n.Doc != nil {
			Walk(v, n.Doc, AddHistory(history, "Doc"), false)
		}
		walkIdentList(v, n.Names, history, "Names")
		Walk(v, n.Type, AddHistory(history, "Type"), true)
		if n.Tag != nil {
			Walk(v, n.Tag, AddHistory(history, "Tag"), false)
		}
		if n.Comment != nil {
			Walk(v, n.Comment, AddHistory(history, "Comment"), false)
		}

	case *ast.FieldList:
		for i, f := range n.List {
			Walk(v, f, AddListHistory(history, "List", i), false)
		}

	// Expressions
	case *ast.BadExpr, *ast.Ident, *ast.BasicLit:
	// nothing to do

	case *ast.Ellipsis:
		if n.Elt != nil {
			Walk(v, n.Elt, AddHistory(history, "Elt"), true)
		}

	case *ast.FuncLit:
		Walk(v, n.Type, AddHistory(history, "Type"), false)
		Walk(v, n.Body, AddHistory(history, "Body"), false)

	case *ast.CompositeLit:
		if n.Type != nil {
			Walk(v, n.Type, AddHistory(history, "Type"), true)
		}
		walkExprList(v, n.Elts, history, "Elts")

	case *ast.ParenExpr:
		Walk(v, n.X, AddHistory(history, "X"), true)

	case *ast.SelectorExpr:
		Walk(v, n.X, AddHistory(history, "X"), true)
		Walk(v, n.Sel, AddHistory(history, "Sel"), true)

	case *ast.IndexExpr:
		Walk(v, n.X, AddHistory(history, "X"), true)
		Walk(v, n.Index, AddHistory(history, "Index"), true)

	case *ast.SliceExpr:
		Walk(v, n.X, AddHistory(history, "X"), true)
		if n.Low != nil {
			Walk(v, n.Low, AddHistory(history, "Low"), true)
		}
		if n.High != nil {
			Walk(v, n.High, AddHistory(history, "High"), true)
		}
		if n.Max != nil {
			Walk(v, n.Max, AddHistory(history, "Max"), true)
		}

	case *ast.TypeAssertExpr:
		Walk(v, n.X, AddHistory(history, "X"), true)

		if n.Type != nil {
			Walk(v, n.Type, AddHistory(history, "Type"), true)
		}

	case *ast.CallExpr:
		Walk(v, n.Fun, AddHistory(history, "Fun"), true)
		walkExprList(v, n.Args, history, "Args")

	case *ast.StarExpr:
		Walk(v, n.X, AddHistory(history, "X"), true)

	case *ast.UnaryExpr:
		Walk(v, n.X, AddHistory(history, "X"), true)

	case *ast.BinaryExpr:
		Walk(v, n.X, AddHistory(history, "X"), true)
		Walk(v, n.Y, AddHistory(history, "Y"), true)

	case *ast.KeyValueExpr:
		Walk(v, n.Key, AddHistory(history, "Key"), true)
		Walk(v, n.Value, AddHistory(history, "Value"), true)

	// Types
	case *ast.ArrayType:
		if n.Len != nil {
			Walk(v, n.Len, AddHistory(history, "Len"), true)
		}
		Walk(v, n.Elt, AddHistory(history, "Elt"), true)

	case *ast.StructType:
		Walk(v, n.Fields, AddHistory(history, "Fields"), false)

	case *ast.FuncType:
		if n.Params != nil {
			Walk(v, n.Params, AddHistory(history, "Params"), false)
		}
		if n.Results != nil {
			Walk(v, n.Results, AddHistory(history, "Results"), false)
		}

	case *ast.InterfaceType:
		Walk(v, n.Methods, AddHistory(history, "Methods"), false)

	case *ast.MapType:
		Walk(v, n.Key, AddHistory(history, "Key"), true)
		Walk(v, n.Value, AddHistory(history, "Value"), true)

	case *ast.ChanType:
		Walk(v, n.Value, AddHistory(history, "Value"), true)

	// Statements
	case *ast.BadStmt:
	// nothing to do

	case *ast.DeclStmt:
		Walk(v, n.Decl, AddHistory(history, "Decl"), true)

	case *ast.EmptyStmt:
	// nothing to do

	case *ast.LabeledStmt:
		Walk(v, n.Label, AddHistory(history, "Label"), false)
		Walk(v, n.Stmt, AddHistory(history, "Stmt"), true)

	case *ast.ExprStmt:
		Walk(v, n.X, AddHistory(history, "X"), true)

	case *ast.SendStmt:
		Walk(v, n.Chan, AddHistory(history, "Chan"), true)
		Walk(v, n.Value, AddHistory(history, "Value"), true)

	case *ast.IncDecStmt:
		Walk(v, n.X, AddHistory(history, "X"), true)

	case *ast.AssignStmt:
		walkExprList(v, n.Lhs, history, "Lhs")
		walkExprList(v, n.Rhs, history, "Rhs")

	case *ast.GoStmt:
		Walk(v, n.Call, AddHistory(history, "Call"), false)

	case *ast.DeferStmt:
		Walk(v, n.Call, AddHistory(history, "Call"), false)

	case *ast.ReturnStmt:
		walkExprList(v, n.Results, history, "Results")

	case *ast.BranchStmt:
		if n.Label != nil {
			Walk(v, n.Label, AddHistory(history, "Label"), false)
		}

	case *ast.BlockStmt:
		walkStmtList(v, n.List, history, "List")

	case *ast.IfStmt:
		if n.Init != nil {
			Walk(v, n.Init, AddHistory(history, "Init"), true)
		}
		Walk(v, n.Cond, AddHistory(history, "Cond"), true)
		Walk(v, n.Body, AddHistory(history, "Body"), false)
		if n.Else != nil {
			Walk(v, n.Else, AddHistory(history, "Else"), true)
		}

	case *ast.CaseClause:
		walkExprList(v, n.List, history, "List")
		walkStmtList(v, n.Body, history, "Body")

	case *ast.SwitchStmt:
		if n.Init != nil {
			Walk(v, n.Init, AddHistory(history, "Init"), true)
		}
		if n.Tag != nil {
			Walk(v, n.Tag, AddHistory(history, "Tag"), true)
		}
		Walk(v, n.Body, AddHistory(history, "Body"), false)

	case *ast.TypeSwitchStmt:
		if n.Init != nil {
			Walk(v, n.Init, AddHistory(history, "Init"), true)
		}
		Walk(v, n.Assign, AddHistory(history, "Assign"), true)
		Walk(v, n.Body, AddHistory(history, "Body"), false)

	case *ast.CommClause:
		if n.Comm != nil {
			Walk(v, n.Comm, AddHistory(history, "Comm"), true)
		}
		walkStmtList(v, n.Body, history, "Body")

	case *ast.SelectStmt:
		Walk(v, n.Body, AddHistory(history, "Body"), false)

	case *ast.ForStmt:
		if n.Init != nil {
			Walk(v, n.Init, AddHistory(history, "Init"), true)
		}
		if n.Cond != nil {
			Walk(v, n.Cond, AddHistory(history, "Cond"), true)
		}
		if n.Post != nil {
			Walk(v, n.Post, AddHistory(history, "Post"), true)
		}
		Walk(v, n.Body, AddHistory(history, "Body"), false)

	case *ast.RangeStmt:
		if n.Key != nil {
			Walk(v, n.Key, AddHistory(history, "Key"), true)
		}
		if n.Value != nil {
			Walk(v, n.Value, AddHistory(history, "Value"), true)
		}
		Walk(v, n.X, AddHistory(history, "X"), true)
		Walk(v, n.Body, AddHistory(history, "Body"), false)

	// Declarations
	case *ast.ImportSpec:
		if n.Doc != nil {
			Walk(v, n.Doc, AddHistory(history, "Doc"), false)
		}
		if n.Name != nil {
			Walk(v, n.Name, AddHistory(history, "Name"), false)
		}
		Walk(v, n.Path, AddHistory(history, "Path"), false)
		if n.Comment != nil {
			Walk(v, n.Comment, AddHistory(history, "Comment"), false)
		}

	case *ast.ValueSpec:
		if n.Doc != nil {
			Walk(v, n.Doc, AddHistory(history, "Doc"), false)
		}
		walkIdentList(v, n.Names, history, "Names")
		if n.Type != nil {
			Walk(v, n.Type, AddHistory(history, "Type"), true)
		}
		walkExprList(v, n.Values, history, "Values")
		if n.Comment != nil {
			Walk(v, n.Comment, AddHistory(history, "Comment"), false)
		}

	case *ast.TypeSpec:
		if n.Doc != nil {
			Walk(v, n.Doc, AddHistory(history, "Doc"), false)
		}
		Walk(v, n.Name, AddHistory(history, "Name"), false)
		Walk(v, n.Type, AddHistory(history, "Type"), true)
		if n.Comment != nil {
			Walk(v, n.Comment, AddHistory(history, "Comment"), false)
		}

	case *ast.BadDecl:
	// nothing to do

	case *ast.GenDecl:
		if n.Doc != nil {
			Walk(v, n.Doc, AddHistory(history, "Doc"), false)
		}
		for i, s := range n.Specs {
			Walk(v, s, AddListHistory(history, "Specs", i), true)
		}

	case *ast.FuncDecl:
		if n.Doc != nil {
			Walk(v, n.Doc, AddHistory(history, "Doc"), false)
		}
		if n.Recv != nil {
			Walk(v, n.Recv, AddHistory(history, "Recv"), false)
		}
		Walk(v, n.Name, AddHistory(history, "Name"), false)
		Walk(v, n.Type, AddHistory(history, "Type"), false)
		if n.Body != nil {
			Walk(v, n.Body, AddHistory(history, "Body"), false)
		}

	// Files and packages
	case *ast.File:
		if n.Doc != nil {
			Walk(v, n.Doc, AddHistory(history, "Doc"), false)
		}
		Walk(v, n.Name, AddHistory(history, "Name"), false)
		walkDeclList(v, n.Decls, history, "Decls")

	case *ast.Package:
		for name, f := range n.Files {
			Walk(v, f, AddHistory(history, fmt.Sprintf("Files[%s]", name)), false)
		}

	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", n))
	}

	v.Visit(nil, nil)
}

type inspector func(ast.Node, []string) bool

func (f inspector) Visit(node ast.Node, history []string) Visitor {
	if f(node, history) {
		return f
	}
	return nil
}

func Inspect(node ast.Node, f func(ast.Node, []string) bool) {
	Walk(inspector(f), node, []string{}, false)
}
