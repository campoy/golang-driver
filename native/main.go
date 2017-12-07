package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"reflect"

	"github.com/sirupsen/logrus"
)

func main() {
	out := json.NewEncoder(os.Stdout)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		input := s.Bytes()
		logrus.Infof("raw request: %s", input)

		var req request
		if err := json.Unmarshal(input, &req); err != nil {
			logrus.Warningf("could not decode request: %v", err)
			continue
		}

		if err := out.Encode(handle(&req)); err != nil {
			logrus.Errorf("could not encode response: %v", err)
		}
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
}

type request struct {
	Content  string
	Language string
}

type response struct {
	Status string
	Errors []string
	AST    interface{}
}

func handle(req *request) *response {
	f, err := parse(req.Content)
	res := &response{
		Status: "ok",
		AST:    struct{ Root interface{} }{f},
	}
	if err != nil {
		res.Status = "fatal"
		res.Errors = append(res.Errors, err.Error())
	}
	return res
}

type node struct {
	InternalType string            `json:",omitempty"`
	InternalName string            `json:",omitempty"`
	Properties   map[string]string `json:",omitempty"`
	Children     []*node           `json:",omitempty"`
	StartOffset  token.Pos         `json:",omitempty"`
	EndOffset    token.Pos         `json:",omitempty"`
}

func parse(content string) (*node, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "", content, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	return tree(fs, f), nil
}

func tree(fs *token.FileSet, n ast.Node) *node {
	v := reflect.ValueOf(n)
	if v.IsNil() {
		return nil
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	root := &node{
		InternalType: t.Name(),
		StartOffset:  n.Pos() - 1,
		EndOffset:    n.End() - 1,
		Properties:   make(map[string]string),
	}

	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		if ignoredFields[name] {
			continue
		}
		field := v.Field(i)
		value := field.Interface()

		switch v := value.(type) {
		case ast.Node:
			if child := tree(fs, v.(ast.Node)); child != nil {
				child.InternalName = name
				root.Children = append(root.Children, child)
			}
			continue
		case nil, token.Pos:
			continue
		default:
		}

		switch field.Kind() {
		case reflect.Slice:
			if field.Len() == 0 {
				continue
			}
			slice := &node{InternalName: name, InternalType: listTypeName(field)}
			for i := 0; i < field.Len(); i++ {
				e := field.Index(i).Interface()
				if n, ok := e.(ast.Node); ok {
					slice.Children = append(slice.Children, tree(fs, n))
				} else {
					panic(fmt.Sprintf("found slice of non nodes: %T", e))
				}
			}
			root.Children = append(root.Children, slice)

		default:
			root.Properties[name] = fmt.Sprint(field)
		}
	}

	return root
}

func listTypeName(v reflect.Value) string {
	t := v.Type().Elem()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return "ListOf" + t.Name()
}

var ignoredFields = map[string]bool{
	"Imports":    true,
	"Scope":      true,
	"Obj":        true,
	"Unresolved": true,
}
