package utils

import "golang.org/x/tools/go/analysis"

// MustNew call given function that creates an analyzer and panic if returned error is not nil
func MustNew(fn func() (*analysis.Analyzer, error)) *analysis.Analyzer {
	a, err := fn()
	if err != nil {
		panic(err)
	}
	return a
}
