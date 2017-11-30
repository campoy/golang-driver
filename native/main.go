package main

import (
	"bufio"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"

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
		AST:    f,
	}
	if err != nil {
		res.Status = "fatal"
		res.Errors = append(res.Errors, err.Error())
	}
	return res
}

func parse(content string) (interface{}, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "", content, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	ast.Walk(astFilter{}, f)
	return struct {
		File *ast.File `json:"file"`
	}{f}, nil
}

// astFilter removes all unnecessary elements for the UAST.
type astFilter struct{}

func (v astFilter) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.Ident:
		n.Obj = nil
	case *ast.File:
		// Imports are already included in Decls.
		n.Imports = nil
		n.Scope = nil
		n.Unresolved = nil
	}
	return v
}
