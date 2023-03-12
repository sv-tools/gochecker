// Package unused integrates the [unused](https://github.com/dominikh/go-tools/tree/master/unused) linter
package unused

import (
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "unused",
	Doc:  "Finds unused code.",
}
