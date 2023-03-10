package skipgenerated

import (
	"fmt"
	"go/ast"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	codeGenerated          = "code generated"
	doNotEdit              = "do not edit"
	autoGeneratedFile      = "autogenerated file"      // easyjson
	automaticallyGenerated = "automatically generated" // genny
)

var (
	GeneratedPhrases = []string{codeGenerated, doNotEdit, automaticallyGenerated, autoGeneratedFile}
	GeneratedRE      = regexp.MustCompile(fmt.Sprintf(`(?i)(%s)`, strings.Join(GeneratedPhrases, "|")))

	Analyzer = &analysis.Analyzer{
		Name:             "skipgenerated",
		Doc:              "filters out the generated files",
		Run:              run,
		RunDespiteErrors: true,
		ResultType:       reflect.TypeOf([]*ast.File{}),
	}
)

func run(pass *analysis.Pass) (any, error) {
	files := make([]*ast.File, 0, len(pass.Files))
	for _, f := range pass.Files {
		if len(f.Comments) > 0 && GeneratedRE.MatchString(f.Comments[0].Text()) {
			continue
		}
		files = append(files, f)
	}
	return files, nil
}
