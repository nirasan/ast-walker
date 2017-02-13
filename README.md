# ast-walker

* ast-walker is golang ast walk helper tool.

# Install

```
go get github.com/nirasan/ast-walker
```

# Usage

## Show ast.Print

* Show whole ast.Print result.

```
ast-walker FILENAME
```

## Show command result

* Show accesser command specified node's ast.Print.
* If you specified 'Decls[0]' command, you get ast.Print of first declaration of root node.

```
ast-walker -command COMMAND FILENAME
```

## Show command matched result

* Show accesser command matched node's ast.Print.
* If you specified '\(*ast.(If|For)Stmt\).Body$' pattern, you get ast.Print of node of 'if' and 'for' statement's body.

```
ast-walker -regex REGEX FILENAME
```

## Show type mached result

* Show type mached node's ast.Print.
* If you specified '*ast.TypeSpec', you get ast.Print of node of type declarations.

```
ast-walker -type NAME FILENAME
```

# Example

## Target File

```
package main

type ST1 struct {
	N int
	B bool
}

type ST2 struct {
	S  string
	IS []int
}
```

## Search first struct declaration by command

```
$ ast-walker -command Decls[0].(*ast.GenDecl).Specs[0] file.go

Found!

COMMAND:
    Decls[0].(*ast.GenDecl).Specs[0]
AST_PRINT:
     0  *ast.TypeSpec {
     1  .  Name: *ast.Ident {
     2  .  .  NamePos: test/file2.go:3:6
     3  .  .  Name: "ST1"
     4  .  .  Obj: *ast.Object {
     5  .  .  .  Kind: type
     6  .  .  .  Name: "ST1"
     7  .  .  .  Decl: *(obj @ 0)
     8  .  .  }
     9  .  }
    10  .  Type: *ast.StructType {
    11  .  .  Struct: test/file2.go:3:10
    12  .  .  Fields: *ast.FieldList {
    13  .  .  .  Opening: test/file2.go:3:17
    14  .  .  .  List: []*ast.Field (len = 2) {
    15  .  .  .  .  0: *ast.Field {
    16  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    17  .  .  .  .  .  .  0: *ast.Ident {
    18  .  .  .  .  .  .  .  NamePos: test/file2.go:4:2
    19  .  .  .  .  .  .  .  Name: "N"
    20  .  .  .  .  .  .  .  Obj: *ast.Object {
    21  .  .  .  .  .  .  .  .  Kind: var
    22  .  .  .  .  .  .  .  .  Name: "N"
    23  .  .  .  .  .  .  .  .  Decl: *(obj @ 15)
    24  .  .  .  .  .  .  .  }
    25  .  .  .  .  .  .  }
    26  .  .  .  .  .  }
    27  .  .  .  .  .  Type: *ast.Ident {
    28  .  .  .  .  .  .  NamePos: test/file2.go:4:4
    29  .  .  .  .  .  .  Name: "int"
    30  .  .  .  .  .  }
    31  .  .  .  .  }
    32  .  .  .  .  1: *ast.Field {
    33  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    34  .  .  .  .  .  .  0: *ast.Ident {
    35  .  .  .  .  .  .  .  NamePos: test/file2.go:5:2
    36  .  .  .  .  .  .  .  Name: "B"
    37  .  .  .  .  .  .  .  Obj: *ast.Object {
    38  .  .  .  .  .  .  .  .  Kind: var
    39  .  .  .  .  .  .  .  .  Name: "B"
    40  .  .  .  .  .  .  .  .  Decl: *(obj @ 32)
    41  .  .  .  .  .  .  .  }
    42  .  .  .  .  .  .  }
    43  .  .  .  .  .  }
    44  .  .  .  .  .  Type: *ast.Ident {
    45  .  .  .  .  .  .  NamePos: test/file2.go:5:4
    46  .  .  .  .  .  .  Name: "bool"
    47  .  .  .  .  .  }
    48  .  .  .  .  }
    49  .  .  .  }
    50  .  .  .  Closing: test/file2.go:6:1
    51  .  .  }
    52  .  .  Incomplete: false
    53  .  }
    54  }
```

## Search struct field declaration by regex

```
$ ast-walker -regex 'Fields.List\[\d+\]$' file.go

Found!

COMMAND:
    Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[0]
AST_PRINT:
     0  *ast.Field {
     1  .  Names: []*ast.Ident (len = 1) {
     2  .  .  0: *ast.Ident {
     3  .  .  .  NamePos: test/file2.go:4:2
     4  .  .  .  Name: "N"
     5  .  .  .  Obj: *ast.Object {
     6  .  .  .  .  Kind: var
     7  .  .  .  .  Name: "N"
     8  .  .  .  .  Decl: *(obj @ 0)
     9  .  .  .  }
    10  .  .  }
    11  .  }
    12  .  Type: *ast.Ident {
    13  .  .  NamePos: test/file2.go:4:4
    14  .  .  Name: "int"
    15  .  }
    16  }

COMMAND:
    Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[1]
AST_PRINT:
     0  *ast.Field {
     1  .  Names: []*ast.Ident (len = 1) {
     2  .  .  0: *ast.Ident {
     3  .  .  .  NamePos: test/file2.go:5:2
     4  .  .  .  Name: "B"
     5  .  .  .  Obj: *ast.Object {
     6  .  .  .  .  Kind: var
     7  .  .  .  .  Name: "B"
     8  .  .  .  .  Decl: *(obj @ 0)
     9  .  .  .  }
    10  .  .  }
    11  .  }
    12  .  Type: *ast.Ident {
    13  .  .  NamePos: test/file2.go:5:4
    14  .  .  Name: "bool"
    15  .  }
    16  }

 COMMAND:
     Decls[1].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[0]
 AST_PRINT:
      0  *ast.Field {
      1  .  Names: []*ast.Ident (len = 1) {
      2  .  .  0: *ast.Ident {
      3  .  .  .  NamePos: test/file2.go:9:2
      4  .  .  .  Name: "S"
      5  .  .  .  Obj: *ast.Object {
      6  .  .  .  .  Kind: var
      7  .  .  .  .  Name: "S"
      8  .  .  .  .  Decl: *(obj @ 0)
      9  .  .  .  }
     10  .  .  }
     11  .  }
     12  .  Type: *ast.Ident {
     13  .  .  NamePos: test/file2.go:9:5
     14  .  .  Name: "string"
     15  .  }
     16  }

 COMMAND:
     Decls[1].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List[1]
 AST_PRINT:
      0  *ast.Field {
      1  .  Names: []*ast.Ident (len = 1) {
      2  .  .  0: *ast.Ident {
      3  .  .  .  NamePos: test/file2.go:10:2
      4  .  .  .  Name: "IS"
      5  .  .  .  Obj: *ast.Object {
      6  .  .  .  .  Kind: var
      7  .  .  .  .  Name: "IS"
      8  .  .  .  .  Decl: *(obj @ 0)
      9  .  .  .  }
     10  .  .  }
     11  .  }
     12  .  Type: *ast.ArrayType {
     13  .  .  Lbrack: test/file2.go:10:5
     14  .  .  Elt: *ast.Ident {
     15  .  .  .  NamePos: test/file2.go:10:7
     16  .  .  .  Name: "int"
     17  .  .  }
     18  .  }
     19  }
```

### Search struct declaration by type

```
$ ast-walker -type *ast.TypeSpec file.go

Found!

COMMAND:
    Decls[0].(*ast.GenDecl).Specs[0]
AST_PRINT:
     0  *ast.TypeSpec {
     1  .  Name: *ast.Ident {
     2  .  .  NamePos: test/file2.go:3:6
     3  .  .  Name: "ST1"
     4  .  .  Obj: *ast.Object {
     5  .  .  .  Kind: type
     6  .  .  .  Name: "ST1"
     7  .  .  .  Decl: *(obj @ 0)
     8  .  .  }
     9  .  }
    10  .  Type: *ast.StructType {
    11  .  .  Struct: test/file2.go:3:10
    12  .  .  Fields: *ast.FieldList {
    13  .  .  .  Opening: test/file2.go:3:17
    14  .  .  .  List: []*ast.Field (len = 2) {
    15  .  .  .  .  0: *ast.Field {
    16  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    17  .  .  .  .  .  .  0: *ast.Ident {
    18  .  .  .  .  .  .  .  NamePos: test/file2.go:4:2
    19  .  .  .  .  .  .  .  Name: "N"
    20  .  .  .  .  .  .  .  Obj: *ast.Object {
    21  .  .  .  .  .  .  .  .  Kind: var
    22  .  .  .  .  .  .  .  .  Name: "N"
    23  .  .  .  .  .  .  .  .  Decl: *(obj @ 15)
    24  .  .  .  .  .  .  .  }
    25  .  .  .  .  .  .  }
    26  .  .  .  .  .  }
    27  .  .  .  .  .  Type: *ast.Ident {
    28  .  .  .  .  .  .  NamePos: test/file2.go:4:4
    29  .  .  .  .  .  .  Name: "int"
    30  .  .  .  .  .  }
    31  .  .  .  .  }
    32  .  .  .  .  1: *ast.Field {
    33  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    34  .  .  .  .  .  .  0: *ast.Ident {
    35  .  .  .  .  .  .  .  NamePos: test/file2.go:5:2
    36  .  .  .  .  .  .  .  Name: "B"
    37  .  .  .  .  .  .  .  Obj: *ast.Object {
    38  .  .  .  .  .  .  .  .  Kind: var
    39  .  .  .  .  .  .  .  .  Name: "B"
    40  .  .  .  .  .  .  .  .  Decl: *(obj @ 32)
    41  .  .  .  .  .  .  .  }
    42  .  .  .  .  .  .  }
    43  .  .  .  .  .  }
    44  .  .  .  .  .  Type: *ast.Ident {
    45  .  .  .  .  .  .  NamePos: test/file2.go:5:4
    46  .  .  .  .  .  .  Name: "bool"
    47  .  .  .  .  .  }
    48  .  .  .  .  }
    49  .  .  .  }
    50  .  .  .  Closing: test/file2.go:6:1
    51  .  .  }
    52  .  .  Incomplete: false
    53  .  }
    54  }

COMMAND:
    Decls[1].(*ast.GenDecl).Specs[0]
AST_PRINT:
     0  *ast.TypeSpec {
     1  .  Name: *ast.Ident {
     2  .  .  NamePos: test/file2.go:8:6
     3  .  .  Name: "ST2"
     4  .  .  Obj: *ast.Object {
     5  .  .  .  Kind: type
     6  .  .  .  Name: "ST2"
     7  .  .  .  Decl: *(obj @ 0)
     8  .  .  }
     9  .  }
    10  .  Type: *ast.StructType {
    11  .  .  Struct: test/file2.go:8:10
    12  .  .  Fields: *ast.FieldList {
    13  .  .  .  Opening: test/file2.go:8:17
    14  .  .  .  List: []*ast.Field (len = 2) {
    15  .  .  .  .  0: *ast.Field {
    16  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    17  .  .  .  .  .  .  0: *ast.Ident {
    18  .  .  .  .  .  .  .  NamePos: test/file2.go:9:2
    19  .  .  .  .  .  .  .  Name: "S"
    20  .  .  .  .  .  .  .  Obj: *ast.Object {
    21  .  .  .  .  .  .  .  .  Kind: var
    22  .  .  .  .  .  .  .  .  Name: "S"
    23  .  .  .  .  .  .  .  .  Decl: *(obj @ 15)
    24  .  .  .  .  .  .  .  }
    25  .  .  .  .  .  .  }
    26  .  .  .  .  .  }
    27  .  .  .  .  .  Type: *ast.Ident {
    28  .  .  .  .  .  .  NamePos: test/file2.go:9:5
    29  .  .  .  .  .  .  Name: "string"
    30  .  .  .  .  .  }
    31  .  .  .  .  }
    32  .  .  .  .  1: *ast.Field {
    33  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    34  .  .  .  .  .  .  0: *ast.Ident {
    35  .  .  .  .  .  .  .  NamePos: test/file2.go:10:2
    36  .  .  .  .  .  .  .  Name: "IS"
    37  .  .  .  .  .  .  .  Obj: *ast.Object {
    38  .  .  .  .  .  .  .  .  Kind: var
    39  .  .  .  .  .  .  .  .  Name: "IS"
    40  .  .  .  .  .  .  .  .  Decl: *(obj @ 32)
    41  .  .  .  .  .  .  .  }
    42  .  .  .  .  .  .  }
    43  .  .  .  .  .  }
    44  .  .  .  .  .  Type: *ast.ArrayType {
    45  .  .  .  .  .  .  Lbrack: test/file2.go:10:5
    46  .  .  .  .  .  .  Elt: *ast.Ident {
    47  .  .  .  .  .  .  .  NamePos: test/file2.go:10:7
    48  .  .  .  .  .  .  .  Name: "int"
    49  .  .  .  .  .  .  }
    50  .  .  .  .  .  }
    51  .  .  .  .  }
    52  .  .  .  }
    53  .  .  .  Closing: test/file2.go:11:1
    54  .  .  }
    55  .  .  Incomplete: false
    56  .  }
    57  }
```
