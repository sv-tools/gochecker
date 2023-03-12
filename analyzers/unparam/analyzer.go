// Package unparam integrates the [unparam](https://github.com/mvdan/unparam) linter
package unparam

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/packages"
	"mvdan.cc/unparam/check"

	"github.com/sv-tools/gochecker/analyzers/skipgenerated"
)

var Analyzer = &analysis.Analyzer{
	Name: "unparam",
	Doc:  "Reports unused function parameters and results in your code.",
	Requires: []*analysis.Analyzer{
		skipgenerated.Analyzer,
		buildssa.Analyzer,
	},
	Run: run,
}

var checkExported bool

func init() {
	Analyzer.Flags.BoolVar(&checkExported, "exported", false, "inspect exported functions")
}

func run(pass *analysis.Pass) (any, error) {
	files := pass.ResultOf[skipgenerated.Analyzer].([]*ast.File)
	ssa := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)

	checker := check.Checker{}
	checker.CheckExportedFuncs(checkExported)
	checker.Packages([]*packages.Package{{
		Fset:      pass.Fset,
		Syntax:    files,
		Types:     pass.Pkg,
		TypesInfo: pass.TypesInfo,
	}})
	checker.ProgramSSA(ssa.Pkg.Prog)

	issues, err := checker.Check()
	if err != nil {
		return nil, err
	}

	for _, issue := range issues {
		pass.Report(analysis.Diagnostic{
			Pos:     issue.Pos(),
			Message: issue.Message(),
		})
	}

	return nil, nil
}
