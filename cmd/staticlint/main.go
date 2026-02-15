// Линтер, объединяющий в себе базовые анализаторы go/analysis,
// публичные анализаторы bodyclose и errcheck,
// анализаторы из staticcheck и кастомный анализатор ExitCheck.
package main

import (
	"go/ast"

	"github.com/kisielk/errcheck/errcheck"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

// Анализатор, проверяющий, что в методе main пакета main нет прямого вызова os.Exit
var ExitCheck = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "checking for of os.exit() in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "main" {
				return true
			}
			ast.Inspect(fn, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				obj, ok := pass.TypesInfo.Uses[sel.Sel]
				if !ok {
					return true
				}
				if obj.Pkg() != nil && obj.Pkg().Path() == "os" && obj.Name() == "Exit" {
					pass.Reportf(call.Pos(), "direct call of os.Exit in main func should be removed")
				}
				return true
			})
			return true
		})
	}
	return nil, nil
}

func main() {
	mychecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		errcheck.Analyzer,
		bodyclose.Analyzer,
		ExitCheck,
	}

	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	multichecker.Main(
		mychecks...,
	)
}
