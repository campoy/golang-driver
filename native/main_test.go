package main

import (
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/bblfsh/sdk.v1/uast"
)

func roles(rs ...uast.Role) []uast.Role { return rs }

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

func node(name string, role uast.Role, props map[string]string, children ...*uast.Node) *uast.Node {
	var roles []uast.Role
	if role != 0 {
		roles = []uast.Role{role}
	}
	return &uast.Node{
		InternalType: name,
		Properties:   props,
		Roles:        roles,
		Children:     children,
	}
}

func TestHandle(t *testing.T) {
	tt := []struct {
		name    string
		content string
		err     string
		ast     *uast.Node
	}{
		{
			name:    "empty file",
			content: "",
			err:     "[1:1: expected 'package', found 'EOF'",
		},
		{
			name:    "just package main",
			content: "package main",
			ast: node("File", uast.File, nil,
				node("Name", uast.Identifier, props("Name", "main")),
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
			ast: node("File", uast.File, nil,
				node("Name", uast.Identifier, props("Name", "main")),
				node("Decls", 0, nil,
					node("GenDecl", 0, props("Tok", "import"),
						node("Specs", 0, nil,
							node("ImportSpec", uast.Import, nil,
								node("Path", 0, props("Kind", "STRING", "Value", "\"fmt\"")),
							),
						),
					),
					node("FuncDecl", uast.Function, nil,
						node("Name", uast.Identifier, props("Name", "main")),
						node("Type", uast.Type, nil,
							node("Params", uast.ArgsList, nil),
						),
						node("Body", uast.Block, nil,
							node("List", 0, nil,
								node("ExprStmt", uast.Statement, nil,
									node("X", uast.Call, nil,
										node("Fun", 0, nil,
											node("X", uast.Identifier, props("Name", "fmt")),
											node("Sel", uast.Identifier, props("Name", "Println")),
										),
										node("Args", 0, nil,
											node("BasicLit", 0, props("Kind", "STRING", "Value", "\"hello\"")),
										),
									),
								),
							),
						),
					),
				),
			),
		},
	}

	ignorePos := cmp.Comparer(func(a, b *uast.Position) bool { return true })

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
