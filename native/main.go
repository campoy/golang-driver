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
	"gopkg.in/bblfsh/sdk.v1/uast"
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

func parse(content string) (*uast.Node, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "", content, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	return tree(fs, f), nil
}

func position(pos token.Pos) *uast.Position {
	return &uast.Position{Offset: uint32(pos)}
}

func tree(fs *token.FileSet, node ast.Node) *uast.Node {
	v := reflect.ValueOf(node)
	if v.IsNil() {
		return nil
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	root := &uast.Node{
		InternalType:  t.Name(),
		StartPosition: position(node.Pos()),
		EndPosition:   position(node.End()),
		Properties:    make(map[string]string),
	}
	if role, ok := rolesByName[root.InternalType]; ok {
		root.Roles = append(root.Roles, role)
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
				child.InternalType = name
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
			slice := &uast.Node{InternalType: name}
			if role, ok := rolesByName[name]; ok {
				slice.Roles = append(slice.Roles, role)
			}
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

var typeOfNode = reflect.TypeOf(new(ast.Node)).Elem()

var ignoredFields = map[string]bool{
	"Imports":    true,
	"Scope":      true,
	"Obj":        true,
	"Unresolved": true,
}

var rolesByName = map[string]uast.Role{
	"Ident":      uast.Identifier,
	"BinaryExpr": uast.Binary,
	"UnaryExpr":  uast.Unary,
	// TODO: which one is correct?
	// "ExprStmt":           uast.Expression,
	"ExprStmt":      uast.Statement,
	"File":          uast.File,
	"Package":       uast.Package,
	"DeclStmt":      uast.Declaration,
	"ImportSpec":    uast.Import,
	"FuncDecl":      uast.Function,
	"FuncLit":       uast.Function,
	"FieldList":     uast.ArgsList,
	"IfStmt":        uast.If,
	"SwitchStmt":    uast.Switch,
	"CaseClause":    uast.Case,
	"ForStmt":       uast.For,
	"BlockStmt":     uast.Block,
	"ReturnStmt":    uast.Return,
	"CallExpr":      uast.Call,
	"SliceExpr":     uast.List,
	"TypeSpec":      uast.Type,
	"ArrayType":     uast.Type,
	"ChanType":      uast.Type,
	"FuncType":      uast.Type,
	"InterfaceType": uast.Type,
	"MapType":       uast.Type,
	"StructType":    uast.Type,
	"Type":          uast.Type,
	"AssignStmt":    uast.Assignment,
	"Comment":       uast.Comment,
	"Var":           uast.Variable,
}
