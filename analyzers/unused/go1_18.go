//go:build go1.18

package unused

import (
	"errors"

	"golang.org/x/tools/go/analysis"
)

func init() {
	Analyzer.Requires = nil
	Analyzer.Run = func(pass *analysis.Pass) (any, error) {
		return nil, errors.New("go v1.18 is not supported")
	}
}
