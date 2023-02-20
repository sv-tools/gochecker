package main

import (
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	parseConfig()
	multichecker.Main(analyzers...)
}
