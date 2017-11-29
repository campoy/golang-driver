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
	Content string `json:"content"`
}

type response struct {
	Status string   `json:"status"`
	Errors []string `json:"errors,omitempty"`
	AST    ast.Node `json:"ast,omitempty"`

	// TODO: what should metadata contain?
}

func handle(req *request) *response {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "", req.Content, parser.ParseComments)
	if err != nil {
		return &response{Status: "fatal", Errors: []string{err.Error()}}
	}

	ast.Walk(astFilter{}, f)

	return &response{Status: "ok", AST: f}
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
