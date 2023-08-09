package analyzers

import (
	"4d63.com/gocheckcompilerdirectives/checkcompilerdirectives"
	"4d63.com/gochecknoglobals/checknoglobals"
	"github.com/Abirdcfly/dupword"
	errname "github.com/Antonboom/errname/pkg/analyzer"
	nilnil "github.com/Antonboom/nilnil/pkg/analyzer"
	"github.com/Djarvur/go-err113"
	exhaustruct "github.com/GaijinEntertainment/go-exhaustruct/v2/pkg/analyzer"
	forbidigo "github.com/ashanbrown/forbidigo/pkg/analyzer"
	makezero "github.com/ashanbrown/makezero/pkg/analyzer"
	cyclop "github.com/bkielbasa/cyclop/pkg/analyzer"
	"github.com/blizzy78/varnamelen"
	"github.com/breml/bidichk/pkg/bidichk"
	"github.com/breml/errchkjson"
	ireturn "github.com/butuzov/ireturn/analyzer"
	"github.com/butuzov/mirror"
	"github.com/charithe/durationcheck"
	"github.com/curioswitch/go-reassign"
	nonamedreturns "github.com/firefart/nonamedreturns/analyzer"
	critic "github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/gostaticanalysis/forcetypeassert"
	"github.com/gostaticanalysis/nilerr"
	"github.com/jingyugao/rowserrcheck/passes/rowserr"
	goprintffuncname "github.com/jirfag/go-printf-func-name/pkg/analyzer"
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
	"github.com/sanposhiho/wastedassign/v2"
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
	magicnumbers "github.com/tommy-muehle/go-mnd/v2"
	"github.com/uudashr/gocognit"
	"github.com/yagipy/maintidx"
	"github.com/ykadowak/zerologlint"
	"gitlab.com/bosi/decorder"
	"go.tmz.dev/musttag"
	"golang.org/x/tools/go/analysis"

	"github.com/sv-tools/gochecker/analyzers/gci"
	"github.com/sv-tools/gochecker/analyzers/gofumpt"
	"github.com/sv-tools/gochecker/analyzers/unparam"
	"github.com/sv-tools/gochecker/analyzers/unused"
	"github.com/sv-tools/gochecker/analyzers/utils"
)

// External is the list of all external analyzers (linters)
var External = []*analysis.Analyzer{
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
	exportloopref.Analyzer,                                 // https://github.com/kyoh86/exportloopref
	forbidigo.NewAnalyzer(),                                // https://github.com/ashanbrown/forbidigo
	forcetypeassert.Analyzer,                               // https://github.com/gostaticanalysis/forcetypeassert
	gci.Analyzer,                                           // https://github.com/daixiang0/gci
	ginkgolinter.Analyzer,                                  // https://github.com/nunnatsa/ginkgolinter
	gocognit.Analyzer,                                      // https://github.com/uudashr/gocognit
	gofumpt.Analyzer,                                       // https://github.com/mvdan/gofumpt
	goprintffuncname.Analyzer,                              // https://github.com/jirfag/go-printf-func-name
	grouper.New(),                                          // https://github.com/leonklingele/grouper
	ineffassign.Analyzer,                                   // https://github.com/gordonklaus/ineffassign
	interfacebloat.New(),                                   // https://github.com/sashamelentyev/interfacebloat
	ireturn.NewAnalyzer(),                                  // https://github.com/butuzov/ireturn
	loggercheck.NewAnalyzer(),                              // https://github.com/timonwong/loggercheck
	magicnumbers.Analyzer,                                  // https://github.com/tommy-muehle/go-mnd
	maintidx.Analyzer,                                      // https://github.com/yagipy/maintidx
	makezero.NewAnalyzer(),                                 // https://github.com/ashanbrown/makezero
	mirror.NewAnalyzer(),                                   // https://github.com/butuzov/mirror
	musttag.New(),                                          // https://github.com/junk1tm/musttag
	nilerr.Analyzer,                                        // https://github.com/gostaticanalysis/nilerr
	nilnil.New(),                                           // https://github.com/Antonboom/nilnil
	nlreturn.NewAnalyzer(),                                 // https://github.com/ssgreg/nlreturn
	noctx.Analyzer,                                         // https://github.com/sonatard/noctx
	nonamedreturns.Analyzer,                                // https://github.com/firefart/nonamedreturns
	nosprintfhostport.Analyzer,                             // https://github.com/stbenjam/no-sprintf-host-port
	paralleltest.NewAnalyzer(),                             // https://github.com/kunwardeep/paralleltest
	predeclared.Analyzer,                                   // https://github.com/nishanths/predeclared
	reassign.NewAnalyzer(),                                 // https://github.com/curioswitch/go-reassign
	rowserr.NewAnalyzer(),                                  // https://github.com/jingyugao/rowserrcheck
	sqlclosecheck.NewAnalyzer(),                            // https://github.com/ryanrolds/sqlclosecheck
	tenv.Analyzer,                                          // https://github.com/sivchari/tenv
	testableexamples.NewAnalyzer(),                         // https://github.com/maratori/testableexamples
	testpackage.NewAnalyzer(),                              // https://github.com/maratori/testpackage
	thelper.NewAnalyzer(),                                  // https://github.com/kulti/thelper
	tparallel.Analyzer,                                     // https://github.com/moricho/tparallel
	unparam.Analyzer,                                       // https://github.com/mvdan/unparam
	unused.Analyzer,                                        // https://github.com/dominikh/go-tools/tree/master/unused
	usestdlibvars.New(),                                    // https://github.com/sashamelentyev/usestdlibvars
	varnamelen.NewAnalyzer(),                               // https://github.com/blizzy78/varnamelen
	wastedassign.Analyzer,                                  // https://github.com/sanposhiho/wastedassign
	zerologlint.Analyzer,                                   // https://github.com/ykadowak/zerologlint

	utils.MustNew(func() (*analysis.Analyzer, error) {
		return exhaustruct.NewAnalyzer(nil, nil) // https://github.com/GaijinEntertainment/go-exhaustruct
	}),
}
