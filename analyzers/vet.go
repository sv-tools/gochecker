package analyzers

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
)

const (
	GoVetName        = "govet"
	GoVetExcludeName = "exclude"
	GoVetExtraName   = "extra"
)

var (
	// dummyGoVet is needed to properly setup and parse the cli flags.
	dummyGoVet = &analysis.Analyzer{
		Name: GoVetName,
		Doc: `Enables the set of official go analyzers.

The list of official go vet analyzers to be enabled:

	asmdecl, assign, atomic, bools, buildtag, cgocall, composites, copylocks, 
	httprespons, loopclosure, lostcancel, nilfunc, printf, shift, stdmethods, 
	structtag, tests, unmarshal, unreachable, unsafeptr, unusedresul.

and list of additional official analyzers that will be enabled if the extra flag is set to true:

	atomicalign, deepequalerrors, directive, errorsas, fieldalignment, framepointer, 
	ifaceassert, nilness, reflectvaluecompare, shadow, sigchanyzer, sortslice,
	stringintconv, timeformat, unusedwrite, usesgenerics.
`,
		Run: func(*analysis.Pass) (any, error) { return nil, nil },
	}

	// GoVet the list of official go analyzers (passes)
	// https://pkg.go.dev/golang.org/x/tools/go/analysis/passes
	GoVet = []*analysis.Analyzer{
		dummyGoVet,

		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shift.Analyzer,
		stdmethods.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	}

	GoVetExtra = []*analysis.Analyzer{
		atomicalign.Analyzer,
		deepequalerrors.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		framepointer.Analyzer,
		ifaceassert.Analyzer,
		nilness.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		sigchanyzer.Analyzer,
		sortslice.Analyzer,
		stringintconv.Analyzer,
		timeformat.Analyzer,
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,
	}
)

func init() {
	dummyGoVet.Flags.Bool(GoVetExtraName, false, "enable all go vet extra passes")
	dummyGoVet.Flags.String(GoVetExcludeName, "", "comma separated list of official analyzers to exclude")
}
