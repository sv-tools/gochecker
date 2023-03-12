package gofumpt

import (
	"fmt"
	"go/ast"
	"os"

	"golang.org/x/tools/go/analysis"
	"mvdan.cc/gofumpt/format"

	"github.com/sv-tools/gochecker/analyzers/skipgenerated"
	"github.com/sv-tools/gochecker/analyzers/utils"
)

const (
	Name = "gofumpt"

	LangFlag   = "lang"
	ModuleFlag = "module"
	ExtraFlag  = "extra"
)

var Analyzer = &analysis.Analyzer{
	Name: Name,
	Doc: `A stricter gofmt

Enforce a stricter format than gofmt, while being backwards compatible. 
That is, gofumpt is happy with a subset of the formats that gofmt is happy with.
`,
	Requires: []*analysis.Analyzer{skipgenerated.Analyzer},
	Run:      run,
}

var options format.Options

func init() {
	Analyzer.Flags.StringVar(&options.LangVersion, LangFlag, "", "The Go language version. The current version will be used if not set.")
	Analyzer.Flags.StringVar(&options.ModulePath, ModuleFlag, "", "The Go module path. The path of current module will be used if not set.")
	Analyzer.Flags.BoolVar(&options.ExtraRules, ExtraFlag, false, "Enables extra formatting rules.")
}

func run(pass *analysis.Pass) (any, error) {
	files := pass.ResultOf[skipgenerated.Analyzer].([]*ast.File)
	if len(files) == 0 {
		return nil, nil
	}

	for _, f := range files {
		fileRef := pass.Fset.File(f.Pos())
		data, err := os.ReadFile(fileRef.Name())
		if err != nil {
			return nil, err
		}
		formatted, err := format.Source(data, options)
		if err != nil {
			return nil, err
		}
		fix, err := utils.GetSuggestedFix(fileRef, data, formatted)
		if err != nil {
			return nil, err
		}
		if fix == nil {
			// no difference
			continue
		}
		pass.Report(analysis.Diagnostic{
			Pos:            fix.TextEdits[0].Pos,
			Message:        fmt.Sprintf("The code is not formatted"),
			SuggestedFixes: []analysis.SuggestedFix{*fix},
		})
	}
	return nil, nil
}
