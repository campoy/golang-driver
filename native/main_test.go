package main

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestHandle(t *testing.T) {
	tt := []struct {
		name    string
		content string
		err     string
		ast     ast.Node
	}{
		{
			name:    "empty file",
			content: "",
			err:     "[1:1: expected 'package', found 'EOF'",
		},
		{
			name:    "just package main",
			content: "package main",
			ast: &ast.File{
				Package: 1,
				Name:    &ast.Ident{NamePos: 9, Name: "main"},
			},
		},
		{
			name:    "hello world",
			content: "package main\nimport \"fmt\"\nfunc main() { fmt.Println(\"hello\") }",
			ast: &ast.File{
				Package: 1,
				Name:    &ast.Ident{NamePos: 9, Name: "main"},
				Decls: []ast.Decl{
					&ast.GenDecl{TokPos: 14, Tok: token.IMPORT, Specs: []ast.Spec{
						&ast.ImportSpec{Path: &ast.BasicLit{ValuePos: 21, Kind: token.STRING, Value: `"fmt"`}},
					}},
					&ast.FuncDecl{
						Name: &ast.Ident{NamePos: 32, Name: "main"},
						Type: &ast.FuncType{Func: 27, Params: &ast.FieldList{Opening: 36, Closing: 37}},
						Body: &ast.BlockStmt{
							Lbrace: 39,
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   &ast.Ident{NamePos: 41, Name: "fmt"},
											Sel: &ast.Ident{NamePos: 45, Name: "Println"},
										},
										Lparen: 52,
										Args: []ast.Expr{&ast.BasicLit{
											ValuePos: 53,
											Kind:     token.STRING,
											Value:    `"hello"`,
										}},
										Rparen: 60},
								},
							},
							Rbrace: 62,
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res := handle(&request{Content: tc.content})
			if tc.err == "" && res.Status != "ok" {
				t.Fatalf("unexpected error: %s (%v)", res.Status, res.Errors)
			}
			if tc.err != "" && res.Status == "ok" {
				t.Fatalf("expected error %q; got ok", tc.err)
			}
			if !cmp.Equal(tc.ast, res.AST, cmpopts.EquateEmpty()) {
				t.Fatalf("different ASTs: %s", cmp.Diff(tc.ast, res.AST, cmpopts.EquateEmpty()))
			}
		})
	}
}
