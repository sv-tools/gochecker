package main

import (
	"gitlab.com/bosi/decorder"
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

	"4d63.com/gocheckcompilerdirectives/checkcompilerdirectives"
	"4d63.com/gochecknoglobals/checknoglobals"
	"github.com/Abirdcfly/dupword"
	errname "github.com/Antonboom/errname/pkg/analyzer"
	nilnil "github.com/Antonboom/nilnil/pkg/analyzer"
	"github.com/Djarvur/go-err113"
	exhaustruct "github.com/GaijinEntertainment/go-exhaustruct/pkg/analyzer"
	forbidigo "github.com/ashanbrown/forbidigo/pkg/analyzer"
	makezero "github.com/ashanbrown/makezero/pkg/analyzer"
	cyclop "github.com/bkielbasa/cyclop/pkg/analyzer"
	"github.com/breml/bidichk/pkg/bidichk"
	"github.com/breml/errchkjson"
	ireturn "github.com/butuzov/ireturn/analyzer"
	"github.com/charithe/durationcheck"
	"github.com/curioswitch/go-reassign"
	gci "github.com/daixiang0/gci/pkg/analyzer"
	nonamedreturns "github.com/firefart/nonamedreturns/analyzer"
	critic "github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/gostaticanalysis/forcetypeassert"
	"github.com/gostaticanalysis/nilerr"
	"github.com/jingyugao/rowserrcheck/passes/rowserr"
	goprintffuncname "github.com/jirfag/go-printf-func-name/pkg/analyzer"
	"github.com/julz/importas"
	"github.com/junk1tm/musttag"
	"github.com/kisielk/errcheck/errcheck"
	"github.com/kkHAIKE/contextcheck"
	thelper "github.com/kulti/thelper/pkg/analyzer"
	"github.com/kunwardeep/paralleltest/pkg/paralleltest"
	"github.com/kyoh86/exportloopref"
	grouper "github.com/leonklingele/grouper/pkg/analyzer"
	"github.com/lufeee/execinquery"
	"github.com/maratori/testableexamples/pkg/testableexamples"
	"github.com/maratori/testpackage/pkg/testpackage"
	"github.com/moricho/tparallel"
	"github.com/nishanths/exhaustive"
	"github.com/nishanths/predeclared/passes/predeclared"
	"github.com/nunnatsa/ginkgolinter"
	"github.com/polyfloyd/go-errorlint/errorlint"
	sqlclosecheck "github.com/ryanrolds/sqlclosecheck/pkg/analyzer"
	interfacebloat "github.com/sashamelentyev/interfacebloat/pkg/analyzer"
	usestdlibvars "github.com/sashamelentyev/usestdlibvars/pkg/analyzer"
	"github.com/sivchari/containedctx"
	"github.com/sivchari/tenv"
	"github.com/sonatard/noctx"
	"github.com/ssgreg/nlreturn/v2/pkg/nlreturn"
	nosprintfhostport "github.com/stbenjam/no-sprintf-host-port/pkg/analyzer"
	"github.com/tdakkota/asciicheck"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"github.com/timonwong/loggercheck"
	"github.com/tommy-muehle/go-mnd/v2"
	"github.com/uudashr/gocognit"
	"github.com/yagipy/maintidx"
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
	asciicheck.NewAnalyzer(),                               // https://github.com/tdakkota/asciicheck
	bidichk.NewAnalyzer(),                                  // https://github.com/breml/bidichk
	bodyclose.Analyzer,                                     // https://github.com/timakin/bodyclose
	checkcompilerdirectives.Analyzer(),                     // https://github.com/leighmcculloch/gocheckcompilerdirectives
	checknoglobals.Analyzer(),                              // https://github.com/leighmcculloch/gochecknoglobals
	containedctx.Analyzer,                                  // https://github.com/sivchari/containedctx
	contextcheck.NewAnalyzer(contextcheck.Configuration{}), // https://github.com/kkHAIKE/contextcheck
	critic.Analyzer,                                        // https://github.com/go-critic/go-critic
	cyclop.NewAnalyzer(),                                   // https://github.com/bkielbasa/cyclop
	decorder.Analyzer,                                      // https://gitlab.com/bosi/decorder
	dupword.NewAnalyzer(),                                  // https://github.com/Abirdcfly/dupword
	durationcheck.Analyzer,                                 // https://github.com/charithe/durationcheck
	err113.NewAnalyzer(),                                   // https://github.com/Djarvur/go-err113.git
	errcheck.Analyzer,                                      // https://github.com/kisielk/errcheck
	errchkjson.NewAnalyzer(),                               // https://github.com/breml/errchkjson
	errname.New(),                                          // https://github.com/Antonboom/errname
	errorlint.NewAnalyzer(),                                // https://github.com/polyfloyd/go-errorlint
	execinquery.Analyzer,                                   // https://github.com/1uf3/execinquery
	exhaustive.Analyzer,                                    // https://github.com/nishanths/exhaustive
	exhaustruct.Analyzer,                                   // https://github.com/GaijinEntertainment/go-exhaustruct
	exportloopref.Analyzer,                                 // https://github.com/kyoh86/exportloopref
	forbidigo.NewAnalyzer(),                                // https://github.com/ashanbrown/forbidigo
	forcetypeassert.Analyzer,                               // https://github.com/gostaticanalysis/forcetypeassert
	gci.Analyzer,                                           // https://github.com/daixiang0/gci
	ginkgolinter.Analyzer,                                  // https://github.com/nunnatsa/ginkgolinter
	gocognit.Analyzer,                                      // https://github.com/uudashr/gocognit
	goprintffuncname.Analyzer,                              // https://github.com/jirfag/go-printf-func-name
	grouper.New(),                                          // https://github.com/leonklingele/grouper
	importas.Analyzer,                                      // https://github.com/julz/importas
	ineffassign.Analyzer,                                   // https://github.com/gordonklaus/ineffassign
	interfacebloat.New(),                                   // https://github.com/sashamelentyev/interfacebloat
	ireturn.NewAnalyzer(),                                  // https://github.com/butuzov/ireturn
	loggercheck.NewAnalyzer(),                              // https://github.com/timonwong/loggercheck
	magic_numbers.Analyzer,                                 // https://github.com/tommy-muehle/go-mnd
	maintidx.Analyzer,                                      // https://github.com/yagipy/maintidx
	makezero.NewAnalyzer(),                                 // https://github.com/ashanbrown/makezero
	musttag.New(),                                          // https://github.com/junk1tm/musttag
	nilerr.Analyzer,                                        // https://github.com/gostaticanalysis/nilerr
	nilnil.New(),                                           // https://github.com/Antonboom/nilnil
	nlreturn.NewAnalyzer(),                                 // https://github.com/ssgreg/nlreturn
	noctx.Analyzer,                                         // https://github.com/sonatard/noctx
	nonamedreturns.Analyzer,                                // https://github.com/firefart/nonamedreturns
	nosprintfhostport.Analyzer,                             // https://github.com/stbenjam/no-sprintf-host-port
	paralleltest.Analyzer,                                  // https://github.com/kunwardeep/paralleltest
	predeclared.Analyzer,                                   // https://github.com/nishanths/predeclared
	reassign.NewAnalyzer(),                                 // https://github.com/curioswitch/go-reassign
	rowserr.NewAnalyzer(),                                  // https://github.com/jingyugao/rowserrcheck
	sqlclosecheck.NewAnalyzer(),                            // https://github.com/ryanrolds/sqlclosecheck
	tenv.Analyzer,                                          // https://github.com/sivchari/tenv
	testableexamples.NewAnalyzer(),                         // https://github.com/maratori/testableexamples
	testpackage.NewAnalyzer(),                              // https://github.com/maratori/testpackage
	thelper.NewAnalyzer(),                                  // https://github.com/kulti/thelper
	tparallel.Analyzer,                                     // https://github.com/moricho/tparallel
	usestdlibvars.New(),                                    // https://github.com/sashamelentyev/usestdlibvars
}
