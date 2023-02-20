package main

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

	"github.com/Abirdcfly/dupword"
	errname "github.com/Antonboom/errname/pkg/analyzer"
	cyclop "github.com/bkielbasa/cyclop/pkg/analyzer"
	"github.com/breml/bidichk/pkg/bidichk"
	"github.com/breml/errchkjson"
	"github.com/charithe/durationcheck"
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/kisielk/errcheck/errcheck"
	"github.com/kkHAIKE/contextcheck"
	"github.com/lufeee/execinquery"
	"github.com/nishanths/exhaustive"
	"github.com/polyfloyd/go-errorlint/errorlint"
	"github.com/sivchari/containedctx"
	"github.com/tdakkota/asciicheck"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"gitlab.com/bosi/decorder"
)

var analyzers = []*analysis.Analyzer{
	// official go analyzers
	// https://pkg.go.dev/golang.org/x/tools/go/analysis/passes
	asmdecl.Analyzer,
	assign.Analyzer,
	atomic.Analyzer,
	atomicalign.Analyzer,
	bools.Analyzer,
	buildtag.Analyzer,
	cgocall.Analyzer,
	composite.Analyzer,
	copylock.Analyzer,
	deepequalerrors.Analyzer,
	directive.Analyzer,
	errorsas.Analyzer,
	fieldalignment.Analyzer,
	framepointer.Analyzer,
	httpresponse.Analyzer,
	ifaceassert.Analyzer,
	loopclosure.Analyzer,
	lostcancel.Analyzer,
	nilfunc.Analyzer,
	nilness.Analyzer,
	printf.Analyzer,
	reflectvaluecompare.Analyzer,
	shadow.Analyzer,
	shift.Analyzer,
	sigchanyzer.Analyzer,
	sortslice.Analyzer,
	stdmethods.Analyzer,
	stringintconv.Analyzer,
	structtag.Analyzer,
	testinggoroutine.Analyzer,
	tests.Analyzer,
	timeformat.Analyzer,
	unmarshal.Analyzer,
	unreachable.Analyzer,
	unsafeptr.Analyzer,
	unusedresult.Analyzer,
	unusedwrite.Analyzer,
	usesgenerics.Analyzer,

	// custom analyzers (linters)
	asciicheck.NewAnalyzer(), // https://github.com/tdakkota/asciicheck
	bidichk.NewAnalyzer(),    // https://github.com/breml/bidichk
	bodyclose.Analyzer,       // https://github.com/timakin/bodyclose
	containedctx.Analyzer,    // https://github.com/sivchari/containedctx
	contextcheck.NewAnalyzer(contextcheck.Configuration{}), // https://github.com/kkHAIKE/contextcheck
	cyclop.NewAnalyzer(),     // https://github.com/bkielbasa/cyclop
	decorder.Analyzer,        // https://gitlab.com/bosi/decorder
	dupword.NewAnalyzer(),    // https://github.com/Abirdcfly/dupword
	durationcheck.Analyzer,   // https://github.com/charithe/durationcheck
	errcheck.Analyzer,        // https://github.com/kisielk/errcheck
	errchkjson.NewAnalyzer(), // https://github.com/breml/errchkjson
	ineffassign.Analyzer,     // https://github.com/gordonklaus/ineffassign
	errname.New(),            // https://github.com/Antonboom/errname
	errorlint.NewAnalyzer(),  // https://github.com/polyfloyd/go-errorlint
	execinquery.Analyzer,     // https://github.com/1uf3/execinquery
	exhaustive.Analyzer,      // https://github.com/nishanths/exhaustive
}