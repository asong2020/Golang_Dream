package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	v := visitor{fset: token.NewFileSet()}
	for _, filePath := range os.Args[1:] {
		if filePath == "-" { // to be able to run this like "go run main.go -- input.go"
			continue
		}

		f, err := parser.ParseFile(v.fset, filePath, nil, 0)
		if err != nil {
			log.Fatalf("Failed to parse file %s: %s", filePath, err)
		}
		ast.Walk(&v, f)
	}
}

type visitor struct {
	fset *token.FileSet
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return v
	}

	params := funcDecl.Type.Params.List // get params
	// list is equal of zero that don't need to checker.
	if len(params) == 0 {
		return v
	}

	firstParamType, ok := params[0].Type.(*ast.SelectorExpr)
	if ok && firstParamType.Sel.Name == "Context" {
		return v
	}

	fmt.Printf("%s: %s function first params should be context\n",
		v.fset.Position(node.Pos()), funcDecl.Name.Name)
	return v
}
