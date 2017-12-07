package normalizer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/bblfsh/sdk.v1/protocol"
	"gopkg.in/bblfsh/sdk.v1/uast"
)

func TestTransformers(t *testing.T) {
	tt := []struct {
		name string
		code string
		in   *uast.Node
		out  *uast.Node
	}{
		{
			name: "empty program",
			code: "",
			in:   &uast.Node{},
			out: &uast.Node{
				Roles: []uast.Role{uast.File},
			},
		},
		{
			name: "package main",
			code: "package main",
			in: &uast.Node{
				InternalType: "File",
				Children: []*uast.Node{
					{
						InternalType: "Ident",
						Properties:   map[string]string{"Name": "Name", "internalRole": "Children"},
						Children: []*uast.Node{
							{
								Properties: map[string]string{"Name": "main", "internalRole": "Properties"},
							},
						},
					},
				},
				StartPosition: &uast.Position{Offset: 0},
				EndPosition:   &uast.Position{Offset: 12},
			},
			out: &uast.Node{
				Roles:        []uast.Role{uast.File},
				InternalType: "File",
				Children: []*uast.Node{
					{
						InternalType: "Ident",
						Roles:        []uast.Role{uast.Identifier},
						Properties:   map[string]string{"Name": "Name", "internalRole": "Children"},
						Children: []*uast.Node{
							{
								Properties: map[string]string{"Name": "main", "internalRole": "Properties"},
							},
						},
					},
				},
				StartPosition: &uast.Position{Offset: 0, Line: 1, Col: 1},
				EndPosition:   &uast.Position{Offset: 12, Line: 1, Col: 13},
			},
		},
		{
			name: "binary ops",
			code: "const a = 3 + 5 * 10",
			in: &uast.Node{
				InternalType: "GenDecl",
				Properties: map[string]string{
					"Tok": "const",
				},
				Children: []*uast.Node{{
					InternalType: "ListOfSpec",
					Properties: map[string]string{
						"InternalName": "Specs",
						"internalRole": "Children",
					},
					Children: []*uast.Node{{
						InternalType: "ValueSpec",
						Properties: map[string]string{
							"internalRole": "Children",
						},
						Children: []*uast.Node{{
							InternalType: "ListOfIdent",
							Properties: map[string]string{
								"InternalName": "Names",
								"internalRole": "Children",
							},
							Children: []*uast.Node{{
								InternalType: "Ident",
								Properties: map[string]string{
									"Name":         "a",
									"internalRole": "Children",
								},
								StartPosition: &uast.Position{Offset: 6},
								EndPosition:   &uast.Position{Offset: 7},
							}},
						}, {
							InternalType: "ListOfExpr",
							Properties: map[string]string{
								"InternalName": "Values",
								"internalRole": "Children",
							},
							Children: []*uast.Node{{
								InternalType: "BinaryExpr",
								Properties: map[string]string{
									"Op":           "+",
									"internalRole": "Children",
								},
								Children: []*uast.Node{{
									InternalType: "BasicLit",
									Properties: map[string]string{
										"Kind":         "INT",
										"Value":        "3",
										"internalRole": "Children",
										"InternalName": "X",
									},
									StartPosition: &uast.Position{Offset: 10},
									EndPosition:   &uast.Position{Offset: 11},
								}, {
									InternalType: "BinaryExpr",
									Properties: map[string]string{
										"InternalName": "Y",
										"Op":           "*",
										"internalRole": "Children",
									},
									Children: []*uast.Node{{
										InternalType: "BasicLit",
										Properties: map[string]string{
											"Value":        "5",
											"internalRole": "Children",
											"InternalName": "X",
											"Kind":         "INT",
										},
										StartPosition: &uast.Position{Offset: 14},
										EndPosition:   &uast.Position{Offset: 15},
									}, {
										InternalType: "BasicLit",
										Properties: map[string]string{
											"InternalName": "Y",
											"Kind":         "INT",
											"Value":        "10",
											"internalRole": "Children",
										},
										StartPosition: &uast.Position{Offset: 18},
										EndPosition:   &uast.Position{Offset: 20},
									}},
									StartPosition: &uast.Position{Offset: 14},
									EndPosition:   &uast.Position{Offset: 20},
								}},
								StartPosition: &uast.Position{Offset: 10},
								EndPosition:   &uast.Position{Offset: 20},
							}},
						}},
						StartPosition: &uast.Position{Offset: 6},
						EndPosition:   &uast.Position{Offset: 20},
					}},
				}},
				StartPosition: &uast.Position{Offset: 0},
				EndPosition:   &uast.Position{Offset: 20},
			},
			out: &uast.Node{
				InternalType: "GenDecl",
				Roles:        []uast.Role{uast.File},
				Properties: map[string]string{
					"Tok": "const",
				},
				Children: []*uast.Node{{
					InternalType: "ListOfSpec",
					Properties: map[string]string{
						"InternalName": "Specs",
						"internalRole": "Children",
					},
					Children: []*uast.Node{{
						InternalType: "ValueSpec",
						Properties: map[string]string{
							"internalRole": "Children",
						},
						Children: []*uast.Node{{
							InternalType: "ListOfIdent",
							Properties: map[string]string{
								"InternalName": "Names",
								"internalRole": "Children",
							},
							Children: []*uast.Node{{
								InternalType: "Ident",
								Roles:        []uast.Role{uast.Identifier},
								Properties: map[string]string{
									"Name":         "a",
									"internalRole": "Children",
								},
								StartPosition: &uast.Position{Offset: 6, Line: 1, Col: 7},
								EndPosition:   &uast.Position{Offset: 7, Line: 1, Col: 8},
							}},
						}, {
							InternalType: "ListOfExpr",
							Properties: map[string]string{
								"InternalName": "Values",
								"internalRole": "Children",
							},
							Children: []*uast.Node{{
								InternalType: "BinaryExpr",
								Properties: map[string]string{
									"Op":           "+",
									"internalRole": "Children",
								},
								Children: []*uast.Node{{
									InternalType: "BasicLit",
									Properties: map[string]string{
										"Kind":         "INT",
										"Value":        "3",
										"internalRole": "Children",
										"InternalName": "X",
									},
									StartPosition: &uast.Position{Offset: 10, Line: 1, Col: 11},
									EndPosition:   &uast.Position{Offset: 11, Line: 1, Col: 12},
								}, {
									InternalType: "BinaryExpr",
									Properties: map[string]string{
										"InternalName": "Y",
										"Op":           "*",
										"internalRole": "Children",
									},
									Children: []*uast.Node{{
										InternalType: "BasicLit",
										Properties: map[string]string{
											"Value":        "5",
											"internalRole": "Children",
											"InternalName": "X",
											"Kind":         "INT",
										},
										StartPosition: &uast.Position{Offset: 14, Line: 1, Col: 15},
										EndPosition:   &uast.Position{Offset: 15, Line: 1, Col: 16},
									}, {
										InternalType: "BasicLit",
										Properties: map[string]string{
											"InternalName": "Y",
											"Kind":         "INT",
											"Value":        "10",
											"internalRole": "Children",
										},
										StartPosition: &uast.Position{Offset: 18, Line: 1, Col: 19},
										EndPosition:   &uast.Position{Offset: 20, Line: 1, Col: 21},
									}},
									StartPosition: &uast.Position{Offset: 14, Line: 1, Col: 15},
									EndPosition:   &uast.Position{Offset: 20, Line: 1, Col: 21},
								}},
								StartPosition: &uast.Position{Offset: 10, Line: 1, Col: 11},
								EndPosition:   &uast.Position{Offset: 20, Line: 1, Col: 21},
							}},
						}},
						StartPosition: &uast.Position{Offset: 6, Line: 1, Col: 7},
						EndPosition:   &uast.Position{Offset: 20, Line: 1, Col: 21},
					},
					},
				}},
				StartPosition: &uast.Position{Offset: 0, Line: 1, Col: 1},
				EndPosition:   &uast.Position{Offset: 20, Line: 1, Col: 21},
			},
		},
	}

	opts := []cmp.Option{}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			node := tc.in
			for _, tr := range Transformers {
				if err := tr.Do(tc.code, protocol.UTF8, node); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			t.Logf("NODE: %s\n", print(t, node, "", "\t"))

			if !cmp.Equal(tc.out, node, opts...) {
				t.Fatalf("found difference: %v", cmp.Diff(tc.out, node, opts...))
			}
		})
	}
}
