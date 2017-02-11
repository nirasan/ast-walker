package ast_walker

import (
	"fmt"
	"go/ast"
	"reflect"
)

type Visitor interface {
	Visit(node ast.Node, history *History) (w Visitor)
}

func walkIdentList(v Visitor, list []*ast.Ident, history *History, name string) {
	for i, x := range list {
		Walk(v, x, history.AppendN(name, i), false)
	}
}

func walkExprList(v Visitor, list []ast.Expr, history *History, name string) {
	for i, x := range list {
		Walk(v, x, history.AppendN(name, i), true)
	}
}

func walkStmtList(v Visitor, list []ast.Stmt, history *History, name string) {
	for i, x := range list {
		Walk(v, x, history.AppendN(name, i), true)
	}
}

func walkDeclList(v Visitor, list []ast.Decl, history *History, name string) {
	for i, x := range list {
		Walk(v, x, history.AppendN(name, i), true)
	}
}

func Walk(v Visitor, node ast.Node, history *History, parentIsInterface bool) {
	if v = v.Visit(node, history); v == nil {
		return
	}

	if parentIsInterface {
		history = history.Append(fmt.Sprintf("(%s)", reflect.TypeOf(node)))
	}

	switch n := node.(type) {
	case *ast.Comment:

	case *ast.CommentGroup:
		for i, c := range n.List {
			Walk(v, c, history.AppendN("List", i), false)
		}

	case *ast.Field:
		if n.Doc != nil {
			Walk(v, n.Doc, history.Append("Doc"), false)
		}
		walkIdentList(v, n.Names, history, "Names")
		Walk(v, n.Type, history.Append("Type"), true)
		if n.Tag != nil {
			Walk(v, n.Tag, history.Append("Tag"), false)
		}
		if n.Comment != nil {
			Walk(v, n.Comment, history.Append("Comment"), false)
		}

	case *ast.FieldList:
		for i, f := range n.List {
			Walk(v, f, history.AppendN("List", i), false)
		}

	// Expressions
	case *ast.BadExpr, *ast.Ident, *ast.BasicLit:
	// nothing to do

	case *ast.Ellipsis:
		if n.Elt != nil {
			Walk(v, n.Elt, history.Append("Elt"), true)
		}

	case *ast.FuncLit:
		Walk(v, n.Type, history.Append("Type"), false)
		Walk(v, n.Body, history.Append("Body"), false)

	case *ast.CompositeLit:
		if n.Type != nil {
			Walk(v, n.Type, history.Append("Type"), true)
		}
		walkExprList(v, n.Elts, history, "Elts")

	case *ast.ParenExpr:
		Walk(v, n.X, history.Append("X"), true)

	case *ast.SelectorExpr:
		Walk(v, n.X, history.Append("X"), true)
		Walk(v, n.Sel, history.Append("Sel"), true)

	case *ast.IndexExpr:
		Walk(v, n.X, history.Append("X"), true)
		Walk(v, n.Index, history.Append("Index"), true)

	case *ast.SliceExpr:
		Walk(v, n.X, history.Append("X"), true)
		if n.Low != nil {
			Walk(v, n.Low, history.Append("Low"), true)
		}
		if n.High != nil {
			Walk(v, n.High, history.Append("High"), true)
		}
		if n.Max != nil {
			Walk(v, n.Max, history.Append("Max"), true)
		}

	case *ast.TypeAssertExpr:
		Walk(v, n.X, history.Append("X"), true)

		if n.Type != nil {
			Walk(v, n.Type, history.Append("Type"), true)
		}

	case *ast.CallExpr:
		Walk(v, n.Fun, history.Append("Fun"), true)
		walkExprList(v, n.Args, history, "Args")

	case *ast.StarExpr:
		Walk(v, n.X, history.Append("X"), true)

	case *ast.UnaryExpr:
		Walk(v, n.X, history.Append("X"), true)

	case *ast.BinaryExpr:
		Walk(v, n.X, history.Append("X"), true)
		Walk(v, n.Y, history.Append("Y"), true)

	case *ast.KeyValueExpr:
		Walk(v, n.Key, history.Append("Key"), true)
		Walk(v, n.Value, history.Append("Value"), true)

	// Types
	case *ast.ArrayType:
		if n.Len != nil {
			Walk(v, n.Len, history.Append("Len"), true)
		}
		Walk(v, n.Elt, history.Append("Elt"), true)

	case *ast.StructType:
		Walk(v, n.Fields, history.Append("Fields"), false)

	case *ast.FuncType:
		if n.Params != nil {
			Walk(v, n.Params, history.Append("Params"), false)
		}
		if n.Results != nil {
			Walk(v, n.Results, history.Append("Results"), false)
		}

	case *ast.InterfaceType:
		Walk(v, n.Methods, history.Append("Methods"), false)

	case *ast.MapType:
		Walk(v, n.Key, history.Append("Key"), true)
		Walk(v, n.Value, history.Append("Value"), true)

	case *ast.ChanType:
		Walk(v, n.Value, history.Append("Value"), true)

	// Statements
	case *ast.BadStmt:
	// nothing to do

	case *ast.DeclStmt:
		Walk(v, n.Decl, history.Append("Decl"), true)

	case *ast.EmptyStmt:
	// nothing to do

	case *ast.LabeledStmt:
		Walk(v, n.Label, history.Append("Label"), false)
		Walk(v, n.Stmt, history.Append("Stmt"), true)

	case *ast.ExprStmt:
		Walk(v, n.X, history.Append("X"), true)

	case *ast.SendStmt:
		Walk(v, n.Chan, history.Append("Chan"), true)
		Walk(v, n.Value, history.Append("Value"), true)

	case *ast.IncDecStmt:
		Walk(v, n.X, history.Append("X"), true)

	case *ast.AssignStmt:
		walkExprList(v, n.Lhs, history, "Lhs")
		walkExprList(v, n.Rhs, history, "Rhs")

	case *ast.GoStmt:
		Walk(v, n.Call, history.Append("Call"), false)

	case *ast.DeferStmt:
		Walk(v, n.Call, history.Append("Call"), false)

	case *ast.ReturnStmt:
		walkExprList(v, n.Results, history, "Results")

	case *ast.BranchStmt:
		if n.Label != nil {
			Walk(v, n.Label, history.Append("Label"), false)
		}

	case *ast.BlockStmt:
		walkStmtList(v, n.List, history, "List")

	case *ast.IfStmt:
		if n.Init != nil {
			Walk(v, n.Init, history.Append("Init"), true)
		}
		Walk(v, n.Cond, history.Append("Cond"), true)
		Walk(v, n.Body, history.Append("Body"), false)
		if n.Else != nil {
			Walk(v, n.Else, history.Append("Else"), true)
		}

	case *ast.CaseClause:
		walkExprList(v, n.List, history, "List")
		walkStmtList(v, n.Body, history, "Body")

	case *ast.SwitchStmt:
		if n.Init != nil {
			Walk(v, n.Init, history.Append("Init"), true)
		}
		if n.Tag != nil {
			Walk(v, n.Tag, history.Append("Tag"), true)
		}
		Walk(v, n.Body, history.Append("Body"), false)

	case *ast.TypeSwitchStmt:
		if n.Init != nil {
			Walk(v, n.Init, history.Append("Init"), true)
		}
		Walk(v, n.Assign, history.Append("Assign"), true)
		Walk(v, n.Body, history.Append("Body"), false)

	case *ast.CommClause:
		if n.Comm != nil {
			Walk(v, n.Comm, history.Append("Comm"), true)
		}
		walkStmtList(v, n.Body, history, "Body")

	case *ast.SelectStmt:
		Walk(v, n.Body, history.Append("Body"), false)

	case *ast.ForStmt:
		if n.Init != nil {
			Walk(v, n.Init, history.Append("Init"), true)
		}
		if n.Cond != nil {
			Walk(v, n.Cond, history.Append("Cond"), true)
		}
		if n.Post != nil {
			Walk(v, n.Post, history.Append("Post"), true)
		}
		Walk(v, n.Body, history.Append("Body"), false)

	case *ast.RangeStmt:
		if n.Key != nil {
			Walk(v, n.Key, history.Append("Key"), true)
		}
		if n.Value != nil {
			Walk(v, n.Value, history.Append("Value"), true)
		}
		Walk(v, n.X, history.Append("X"), true)
		Walk(v, n.Body, history.Append("Body"), false)

	// Declarations
	case *ast.ImportSpec:
		if n.Doc != nil {
			Walk(v, n.Doc, history.Append("Doc"), false)
		}
		if n.Name != nil {
			Walk(v, n.Name, history.Append("Name"), false)
		}
		Walk(v, n.Path, history.Append("Path"), false)
		if n.Comment != nil {
			Walk(v, n.Comment, history.Append("Comment"), false)
		}

	case *ast.ValueSpec:
		if n.Doc != nil {
			Walk(v, n.Doc, history.Append("Doc"), false)
		}
		walkIdentList(v, n.Names, history, "Names")
		if n.Type != nil {
			Walk(v, n.Type, history.Append("Type"), true)
		}
		walkExprList(v, n.Values, history, "Values")
		if n.Comment != nil {
			Walk(v, n.Comment, history.Append("Comment"), false)
		}

	case *ast.TypeSpec:
		if n.Doc != nil {
			Walk(v, n.Doc, history.Append("Doc"), false)
		}
		Walk(v, n.Name, history.Append("Name"), false)
		Walk(v, n.Type, history.Append("Type"), true)
		if n.Comment != nil {
			Walk(v, n.Comment, history.Append("Comment"), false)
		}

	case *ast.BadDecl:
	// nothing to do

	case *ast.GenDecl:
		if n.Doc != nil {
			Walk(v, n.Doc, history.Append("Doc"), false)
		}
		for i, s := range n.Specs {
			Walk(v, s, history.AppendN("Specs", i), true)
		}

	case *ast.FuncDecl:
		if n.Doc != nil {
			Walk(v, n.Doc, history.Append("Doc"), false)
		}
		if n.Recv != nil {
			Walk(v, n.Recv, history.Append("Recv"), false)
		}
		Walk(v, n.Name, history.Append("Name"), false)
		Walk(v, n.Type, history.Append("Type"), false)
		if n.Body != nil {
			Walk(v, n.Body, history.Append("Body"), false)
		}

	// Files and packages
	case *ast.File:
		if n.Doc != nil {
			Walk(v, n.Doc, history.Append("Doc"), false)
		}
		Walk(v, n.Name, history.Append("Name"), false)
		walkDeclList(v, n.Decls, history, "Decls")

	case *ast.Package:
		for name, f := range n.Files {
			Walk(v, f, history.Append(fmt.Sprintf("Files[%s]", name)), false)
		}

	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", n))
	}

	v.Visit(nil, nil)
}

type inspector func(ast.Node, *History) bool

func (f inspector) Visit(node ast.Node, history *History) Visitor {
	if f(node, history) {
		return f
	}
	return nil
}

func Inspect(node ast.Node, f func(ast.Node, *History) bool) {
	Walk(inspector(f), node, NewHistory(0), false)
}
