package main

import (
	"go/token"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func props(vs ...string) map[string]string {
	m := make(map[string]string)
	if len(vs)%2 != 0 {
		log.Fatal("bad props value, list is not even")
	}
	for i := 0; i < len(vs); i += 2 {
		m[vs[i]] = vs[i+1]
	}
	return m
}

func n(name, typ string, props map[string]string, children ...*node) *node {
	return &node{
		InternalName: name,
		InternalType: typ,
		Properties:   props,
		Children:     children,
	}
}

func TestHandle(t *testing.T) {
	tt := []struct {
		name    string
		content string
		err     string
		ast     *node
	}{
		{
			name:    "empty file",
			content: "",
			err:     "[1:1: expected 'package', found 'EOF'",
		},
		{
			name:    "just package main",
			content: "package main",
			ast: n("", "File", nil,
				n("Name", "Ident", props("Name", "main")),
			),
		},
		{
			name: "hello world",
			content: `
				package main
				
				import "fmt"
				
				func main() {
					fmt.Println("hello")
				}`,
			ast: n("", "File", nil,
				n("Name", "Ident", props("Name", "main")),
				n("Decls", "ListOfDecl", nil,
					n("", "GenDecl", props("Tok", "import"),
						n("Specs", "ListOfSpec", nil,
							n("", "ImportSpec", nil,
								n("Path", "BasicLit", props("Kind", "STRING", "Value", "\"fmt\"")),
							),
						),
					),
					n("", "FuncDecl", nil,
						n("Name", "Ident", props("Name", "main")),
						n("Type", "FuncType", nil,
							n("Params", "FieldList", nil),
						),
						n("Body", "BlockStmt", nil,
							n("List", "ListOfStmt", nil,
								n("", "ExprStmt", nil,
									n("X", "CallExpr", nil,
										n("Fun", "SelectorExpr", nil,
											n("X", "Ident", props("Name", "fmt")),
											n("Sel", "Ident", props("Name", "Println")),
										),
										n("Args", "ListOfExpr", nil,
											n("", "BasicLit", props("Kind", "STRING", "Value", "\"hello\"")),
										),
									),
								),
							),
						),
					),
				),
			),
		},
		{
			name: "constant definition",
			content: `
				package constants
				
				const a = 40 + 2`,
			ast: n("", "File", nil,
				n("Name", "Ident", props("Name", "constants")),
				n("Decls", "ListOfDecl", nil,
					n("", "GenDecl", props("Tok", "const"),
						n("Specs", "ListOfSpec", nil,
							n("", "ValueSpec", nil,
								n("Names", "ListOfIdent", nil,
									n("", "Ident", props("Name", "a")),
								),
								n("Values", "ListOfExpr", nil,
									n("", "BinaryExpr", props("Op", "+"),
										n("X", "BasicLit", props("Kind", "INT", "Value", "40")),
										n("Y", "BasicLit", props("Kind", "INT", "Value", "2")),
									),
								),
							),
						),
					),
				),
			),
		},
	}

	ignorePos := cmp.Comparer(func(a, b token.Pos) bool { return true })

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parse(tc.content)
			if tc.err == "" && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if tc.err != "" && err == nil {
				t.Fatalf("expected error %q; got ok", tc.err)
			}
			if !cmp.Equal(tc.ast, res, cmpopts.EquateEmpty(), ignorePos) {
				t.Fatalf("different ASTs: %s", cmp.Diff(tc.ast, res, cmpopts.EquateEmpty(), ignorePos))
			}
		})
	}
}
