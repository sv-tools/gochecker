package analyzers

import "golang.org/x/tools/go/analysis"

// Analyzers is the list of all supported analyzers, including govet and the external
var Analyzers []*analysis.Analyzer

func init() {
	Analyzers = append(Analyzers, GoVet...)
	Analyzers = append(Analyzers, External...)
}
