package firstparamcontext

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "firstparamcontext",
	Doc:  "Checks that functions first param type is Context",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := func(node ast.Node) bool {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok {
			return true
		}

		params := funcDecl.Type.Params.List // get params
		// list is equal of zero that don't need to checker.
		if len(params) == 0 {
			return true
		}

		firstParamType, ok := params[0].Type.(*ast.SelectorExpr)
		if ok && firstParamType.Sel.Name == "Context" {
			return true
		}

		pass.Reportf(node.Pos(), "''%s' function first params should be Context\n",
			funcDecl.Name.Name)
		return true
	}

	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}
	return nil, nil
}
