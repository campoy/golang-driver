package normalizer

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"gopkg.in/bblfsh/sdk.v1/uast"
)

type m = map[string]interface{}

func TestToNode(t *testing.T) {
	tt := []struct {
		name string
		in   m
		out  *uast.Node
	}{
		{
			name: "type goes to type identifier",
			in: m{
				"InternalType": "Ident",
				"InternalName": "Name",
			},
			out: &uast.Node{
				InternalType: "Ident",
				Properties:   map[string]string{"InternalName": "Name"},
			},
		},
		{
			name: "start and end offsets",
			in: m{
				"StartOffset": "5",
				"EndOffset":   "10",
			},
			out: &uast.Node{
				StartPosition: &uast.Position{Offset: 5},
				EndPosition:   &uast.Position{Offset: 10},
			},
		},
		{
			name: "other properties",
			in: m{
				"InternalType": "BasicLit",
				"Properties": m{
					"Kind":  "STRING",
					"Value": "\"hello\"",
				},
			},
			out: &uast.Node{
				InternalType: "BasicLit",
				Properties: map[string]string{
					"Kind":  "STRING",
					"Value": "\"hello\"",
				},
			},
		},
		{
			name: "repeated properties",
			in: m{
				"InternalType": "BasicLit",
				"Kind":         "Foo",
				"Properties": m{
					"Kind": "Bar",
				},
			},
			out: &uast.Node{
				InternalType: "BasicLit",
				Properties: map[string]string{
					"Kind": "Foo",
				},
			},
		},
		{
			name: "package main",
			in: m{
				"InternalType": "File",
				"Children": []interface{}{
					m{
						"InternalType": "Ident",
						"InternalName": "Name",
						"Properties": m{
							"Name": "main",
						},
					},
				},
				"StartOffset": "0",
				"EndOffset":   "12",
			},
			out: &uast.Node{
				InternalType: "File",
				Children: []*uast.Node{{
					InternalType: "Ident",
					Properties: map[string]string{
						"InternalName": "Name",
						"Name":         "main",
						"internalRole": "Children",
					},
				}},
				StartPosition: &uast.Position{Offset: 0},
				EndPosition:   &uast.Position{Offset: 12},
			},
		},
		{
			name: "const a = 3 + 5 * 10",
			in: map[string]interface{}{
				"InternalType": "GenDecl",
				"Properties": map[string]interface{}{
					"Tok": "const",
				},
				"Children": []interface{}{
					map[string]interface{}{
						"InternalName": "Specs",
						"InternalType": "ListOfSpec",
						"Children": []interface{}{
							map[string]interface{}{
								"InternalType": "ValueSpec",
								"Children": []interface{}{
									map[string]interface{}{
										"InternalName": "Names",
										"InternalType": "ListOfIdent",
										"Children": []interface{}{
											map[string]interface{}{
												"InternalType": "Ident",
												"Properties": map[string]interface{}{
													"Name": "a",
												},
												"StartOffset": "6",
												"EndOffset":   "7",
											},
										},
									},
									map[string]interface{}{
										"InternalName": "Values",
										"InternalType": "ListOfExpr",
										"Children": []interface{}{
											map[string]interface{}{
												"InternalType": "BinaryExpr",
												"Properties": map[string]interface{}{
													"Op": "+",
												},
												"Children": []interface{}{
													map[string]interface{}{
														"InternalName": "X",
														"InternalType": "BasicLit",
														"Properties": map[string]interface{}{
															"Kind":  "INT",
															"Value": "3",
														},
														"StartOffset": "10",
														"EndOffset":   "11",
													},
													map[string]interface{}{
														"InternalName": "Y",
														"InternalType": "BinaryExpr",
														"Properties": map[string]interface{}{
															"Op": "*",
														},
														"Children": []interface{}{
															map[string]interface{}{
																"InternalName": "X",
																"InternalType": "BasicLit",
																"Properties": map[string]interface{}{
																	"Kind":  "INT",
																	"Value": "5",
																},
																"StartOffset": "14",
																"EndOffset":   "15",
															},
															map[string]interface{}{
																"InternalName": "Y",
																"InternalType": "BasicLit",
																"Properties": map[string]interface{}{
																	"Kind":  "INT",
																	"Value": "10",
																},
																"StartOffset": "18",
																"EndOffset":   "20",
															},
														},
														"StartOffset": "14",
														"EndOffset":   "20",
													},
												},
												"StartOffset": "10",
												"EndOffset":   "20",
											},
										},
									},
								},
								"StartOffset": "6",
								"EndOffset":   "20",
							},
						},
					},
				},
				"StartOffset": "0",
				"EndOffset":   "20",
			},
			out: &uast.Node{
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
		},
	}

	ignoreRoles := cmp.Comparer(func(a, b []uast.Role) bool { return true })
	opts := []cmp.Option{ignoreRoles, cmpopts.EquateEmpty()}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			node, err := ToNode.ToNode(tc.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			t.Logf("NODE: %s\n", print(t, node, "", "\t"))
			if !cmp.Equal(tc.out, node, opts...) {
				t.Fatalf("found difference: %s", cmp.Diff(tc.out, node, opts...))
			}
		})
	}
}

// print prints a node so you can copy and paste it into the out field of a test.
func print(t *testing.T, n *uast.Node, base, indent string) string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "%s&uast.Node{\n", base)
	fmt.Fprintf(w, "%sInternalType: %q,\n", base+indent, n.InternalType)

	if len(n.Roles) != 0 {
		fmt.Fprintf(w, "%sRoles: []uast.Role{", base+indent)
		for _, r := range n.Roles {
			fmt.Fprintf(w, uast.Role_name[int32(r)])
		}
		fmt.Fprintf(w, "},\n")
	}

	if len(n.Properties) != 0 {
		fmt.Fprintf(w, "%sProperties: map[string]string{\n", base+indent)
		for k, v := range n.Properties {
			fmt.Fprintf(w, "%s%q: %q,\n", base+indent+indent, k, v)
		}
		fmt.Fprintf(w, "%s},\n", base+indent)
	}

	if len(n.Children) != 0 {
		fmt.Fprintf(w, "%sChildren: []*uast.Node{\n", base+indent)
		for _, c := range n.Children {
			fmt.Fprint(w, print(t, c, base+indent, indent))
			fmt.Fprintf(w, ",\n")
		}
		fmt.Fprintf(w, "%s},\n", base+indent)
	}

	if n.StartPosition != nil {
		fmt.Fprintf(w, "%sStartPosition: &uast.Position{Offset: %v},\n", base+indent, n.StartPosition.Offset)
	}
	if n.EndPosition != nil {
		fmt.Fprintf(w, "%sEndPosition: &uast.Position{Offset: %v},\n", base+indent, n.EndPosition.Offset)
	}
	fmt.Fprintf(w, "%s}", base)
	return w.String()
}
