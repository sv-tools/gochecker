//go:build go1.19

package unused

import (
	"fmt"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/unused"
)

func init() {
	Analyzer.Run = run
	Analyzer.Requires = []*analysis.Analyzer{unused.Analyzer.Analyzer}
}

func getKey(obj unused.Object) string {
	return fmt.Sprintf("%s %d %s", obj.Position.Filename, obj.Position.Line, obj.Name)
}

func run(pass *analysis.Pass) (any, error) {
	res := pass.ResultOf[unused.Analyzer.Analyzer].(unused.Result)
	used := make(map[string]struct{})
	for _, obj := range res.Used {
		used[getKey(obj)] = struct{}{}
	}
	if len(res.Used) > 0 {
		files := make(map[string]token.Pos, len(pass.Files))
		for _, f := range pass.Files {
			pos := f.Pos()
			position := pass.Fset.Position(pos)
			files[position.Filename] = pos
		}
		for _, obj := range res.Unused {
			if obj.Kind == "type param" {
				continue
			}

			if _, ok := used[getKey(obj)]; ok {
				continue
			}

			pos, ok := files[obj.Position.Filename]
			if !ok {
				return nil, fmt.Errorf("invalid position: %s", obj.Position)
			}
			pass.Report(analysis.Diagnostic{
				Pos:     token.Pos(int(pos) + obj.Position.Offset),
				Message: fmt.Sprintf("%s %s is unused", obj.Kind, obj.Name),
			})
		}
	}
	return nil, nil
}
