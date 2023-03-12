# gochecker

[![Code Analysis](https://github.com/sv-tools/gochecker/actions/workflows/checks.yaml/badge.svg)](https://github.com/sv-tools/gochecker/actions/workflows/checks.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sv-tools/gochecker.svg)](https://pkg.go.dev/github.com/sv-tools/gochecker)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/sv-tools/gochecker?style=flat)](https://github.com/sv-tools/gochecker/releases)

Go tool to check the code (linter)

Another variation of the go linters.
The `gochecker` is a wrapper for [multichecker](https://pkg.go.dev/golang.org/x/tools/go/analysis/multichecker).
The `gochecker` supports `go vet` interface and includes all official [analyzers](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes)
and some custom linters, see [analyzers.go](analyzers.go) file for the full list.
In theory the `gochecker` can support any analyzers which implement [Analyzer type](https://pkg.go.dev/golang.org/x/tools/go/analysis#Analyzer).

The `gochecker` was implemented in order to solve some limitations of the `golangci-linter`.

## Installation

```shell
go install github.com/sv-tools/gochecker@latest
```

## Usage

```shell
gochecker -config config.yaml ./...
```

or using cli flags:

```shell
gochecker -fieldalignment -fix ./...
```

and please check `gochecker help` or `gochecker help <analyzer>` for full help.

### GitHub Action

```yaml
jobs:
  gochecker:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3
      - name: Install Go
        uses: actions/setup-go@v3
      - name: Install gochecker
        run: go install github.com/sv-tools/gochecker@latest
      - name: gochecker
        run: gochecker -config gochecker.yaml -output github ./...
```

### GitHub Action with reusable workflow

```yaml
jobs:
  gochecker:
    uses: sv-tools/gochecker/.github/workflows/gochecker.yaml@main
    with:
      config: gochecker.yaml
      version: latest # optional; `latest` by default
      args: # optional; any additional command-line arguments
      go-version: # optional; the version of go to be used
```

## Supported analyzers

### Go passes

- `govet` is an aggregator includes all passes mentioned here: https://pkg.go.dev/cmd/vet.

Or all the passes can be added as individual analyzers as well, which allows more precisely confixwguration of the `gochecker`.

- [asmdecl](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/asmdecl) reports mismatches between assembly files and Go declarations. 
- [assign](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/assign) detects useless assignments.
- [atomic](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/atomic) checks for common mistakes using the sync/atomic package.
- [atomicalign](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/atomicalign) checks for non-64-bit-aligned arguments to sync/atomic functions. On non-32-bit platforms, those functions panic if their argument variables are not 64-bit aligned. It is therefore the caller's responsibility to arrange for 64-bit alignment of such variables.
- [bools](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/bools) detects common mistakes involving boolean operators.
- [buildtag](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/buildtag) checks build tags.
- [cgocall](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/cgocall) detects some violations of the cgo pointer passing rules.
- [composite](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/composite) checks for unkeyed composite literals.
- [copylock](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/copylock) checks for locks erroneously passed by value.
- [deepequalerrors](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/deepequalerrors) checks for the use of reflect.DeepEqual with error values.
- [directive](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/directive) checks known Go toolchain directives.
- [errorsas](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/errorsas) checks that the second argument to errors.As is a pointer to a type implementing error.
- [fieldalignment](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/fieldalignment) detects structs that would use less memory if their fields were sorted.
- [framepointer](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/framepointer) reports assembly code that clobbers the frame pointer before saving it.
- [httpresponse](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/httpresponse) checks for mistakes using HTTP responses.
- [ifaceassert](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/ifaceassert) flags impossible interface-interface type assertions.
- [loopclosure](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/loopclosure) checks for references to enclosing loop variables from within nested functions.
- [lostcancel](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/lostcancel) checks for failure to call a context cancellation function.
- [nilfunc](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/nilfunc) checks for useless comparisons against nil.
- [nilness](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/nilness) inspects the control-flow graph of an SSA function and reports errors such as nil pointer dereferences and degenerate nil pointer comparisons.
- [printf](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/printf) checks consistency of Printf format strings and arguments.
- [reflectvaluecompare](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/reflectvaluecompare) checks for accidentally using == or reflect.DeepEqual to compare reflect.Value values. See issues 43993 and 18871.
- [shadow](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/shadow) checks for shadowed variables.
- [shift](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/shift) checks for shifts that exceed the width of an integer.
- [sigchanyzer](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/sigchanyzer) detects misuse of unbuffered signal as argument to signal.Notify.
- [sortslice](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/sortslice) checks for calls to sort.Slice that do not use a slice type as first argument.
- [stdmethods](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/stdmethods) checks for misspellings in the signatures of methods similar to well-known interfaces.
- [stringintconv](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/stringintconv) flags type conversions from integers to strings.
- [structtag](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/structtag) checks struct field tags are well-formed.
- [testinggoroutine](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/testinggoroutine) report calls to (*testing.T).Fatal from goroutines started by a test.
- [tests](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/tests) checks for common mistaken usages of tests and examples.
- [timeformat](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/timeformat) checks for the use of time.Format or time.Parse calls with a bad format.
- [unmarshal](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/unmarshal) checks for passing non-pointer or non-interface types to unmarshal and decode functions.
- [unreachable](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/unreachable) checks for unreachable code.
- [unsafeptr](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/unsafeptr) checks for invalid conversions of uintptr to unsafe.Pointer.
- [unusedresult](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/unusedresult) checks for unused results of calls to certain pure functions.
- [unusedwrite](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/unusedwrite) checks for unused writes to the elements of a struct or array object.
- [usesgenerics](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/usesgenerics) checks for usage of generic features added in Go 1.18.

### External linters

The list of analyzers was taken from the `golangci-lint` 
and then each analyzer was checked and imported if it provides an object of the `Analyzer` type. 

- [asciicheck](https://github.com/tdakkota/asciicheck) checks that your code does not contain non-ASCII identifiers.
- [bidichk](https://github.com/breml/bidichk) checks for dangerous unicode character sequences.
- [bodyclose](https://github.com/timakin/bodyclose) checks whether `res.Body` is correctly closed.
- [checkcompilerdirectives](https://github.com/leighmcculloch/gocheckcompilerdirectives) checks that go compiler directives (//go: comments) are valid and catch easy mistakes.
- [checknoglobals](https://github.com/leighmcculloch/gochecknoglobals) check that no globals are present in Go code.
- [containedctx](https://github.com/sivchari/containedctx) detects struct contained context.Context field. This is discouraged technique in favour of passing context as first argument of method or function.
- [contextcheck](https://github.com/kkHAIKE/contextcheck) checks whether the function uses a non-inherited context, which will result in a broken call link.
- [cyclop](https://github.com/bkielbasa/cyclop) calculates cyclomatic complexities of functions or packages in Go source code.
- [dupword](https://github.com/Abirdcfly/dupword) checks for duplicate words in the source code (usually miswritten).
- [durationcheck](https://github.com/charithe/durationcheck) detects cases where two `time.Duration` values are being multiplied in possibly erroneous ways.
- [err113](https://github.com/Djarvur/go-err113) checks the errors handling expressions.
- [errcheck](https://github.com/kisielk/errcheck) checks for unchecked errors in go programs.
- [errchkjson](https://github.com/breml/errchkjson) checks types passed to the json encoding functions. Reports unsupported types and reports occurrences where the check for the returned error can be omitted.
- [errname](https://github.com/Antonboom/errname) checks that sentinel errors are prefixed with the `Err` and error types are suffixed with the `Error`.
- [errorlint](https://github.com/polyfloyd/go-errorlint) finds code that will cause problems with the error wrapping scheme introduced in Go 1.13.
- [execinquery](https://github.com/lufeee/execinquery) is a simple query string checker in Query function.
- [exhaustive](https://github.com/nishanths/exhaustive) checks exhaustiveness of switch statements of enum-like constants in Go source code.
- [exhaustruct](https://github.com/GaijinEntertainment/go-exhaustruct) finds structures with uninitialized fields.
- [exportloopref](https://github.com/kyoh86/exportloopref) finds exporting pointers for loop variables.
- [forbidigo](https://github.com/ashanbrown/forbidigo) forbids usage of particular identifiers.
- [forcetypeassert](https://github.com/gostaticanalysis/forcetypeassert) finds type assertions which did forcely.
- [gci](https://github.com/daixiang0/gci) controls golang package import order and makes it always deterministic.
- [ginkgolinter](https://github.com/nunnatsa/ginkgolinter) enforces some standards while using the ginkgo and gomega packages.
- [gocognit](https://github.com/uudashr/gocognit) calculates cognitive complexities of functions in Go source code. A measurement of how hard does the code is intuitively to understand.
- [gofumpt](https://github.com/mvdan/gofumpt) enforce a stricter format than gofmt, while being backwards compatible.
- [goprintffuncname](https://github.com/jirfag/go-printf-func-name) checks that printf-like functions are named with f at the end.
- [grouper](https://github.com/leonklingele/grouper) analyzes expression groups.
- [ineffassign](https://github.com/gordonklaus/ineffassign) detects ineffectual assignments in Go code. An assignment is ineffectual if the variable assigned is not thereafter used.
- [interfacebloat](https://github.com/sashamelentyev/interfacebloat) checks length of interface.
- [ireturn](https://github.com/butuzov/ireturn) accept interfaces, return concrete types.
- [loggercheck](https://github.com/timonwong/loggercheck) checks the odd number of key and value pairs for common logger libraries.
- [maintidx](https://github.com/yagipy/maintidx) measures the maintainability index of each function.
- [makezero](https://github.com/ashanbrown/makezero) finds slice declarations that are not initialized with zero length and are later used with append.
- [mnd or magic_number](https://github.com/tommy-muehle/go-mnd) detects magic numbers.
- [musttag](https://github.com/junk1tm/musttag) checks that exported fields of a struct passed to a Marshal-like function are annotated with the relevant tag.
- [nilerr](https://github.com/gostaticanalysis/nilerr) finds code which returns nil even though it checks that error is not nil.
- [nilnil](https://github.com/Antonboom/nilnil) checks that there is no simultaneous return of `nil` error and an invalid value.
- [nlreturn](https://github.com/ssgreg/nlreturn) requires a new line before return and branch statements except when the return is alone inside a statement group (such as an if statement) to increase code clarity.
- [noctx](https://github.com/sonatard/noctx) finds sending http request without `context.Context`.
- [nonamedreturns](https://github.com/firefart/nonamedreturns) reports all named returns.
- [nosprintfhostport](https://github.com/stbenjam/no-sprintf-host-port) checks that sprintf is not used to construct a host:port combination in a URL.
- [paralleltest](https://github.com/kunwardeep/paralleltest) checks that the `t.Parallel` gets called for the test method and for the range of test cases within the test.
- [predeclared](https://github.com/nishanths/predeclared) finds code that overrides one of Go's predeclared identifiers (`new`, `make`, `append`, `uint`, etc.).
- [reassign](https://github.com/curioswitch/go-reassign) detects when reassigning a top-level variable in another package.
- [rowserrcheck](https://github.com/jingyugao/rowserrcheck) checks whether sql.Rows.Err is correctly checked.
- [ruleguard or go-critic](https://github.com/go-critic/go-critic) is the most opinionated Go source code linter.
- [sqlclosecheck](https://github.com/ryanrolds/sqlclosecheck) checks if SQL rows/statements are closed. Unclosed rows and statements may cause DB connection pool exhaustion.
- [tenv](https://github.com/sivchari/tenv) detects using os.Setenv instead of `t.Setenv` since Go1.17.
- [testableexamples](https://github.com/maratori/testableexamples)
- [testpackage](https://github.com/maratori/testpackage) checks if examples are testable (have an expected output).
- [thelper](https://github.com/kulti/thelper) detects golang test helpers without `t.Helper()` call. Also, it checks the consistency of test helpers and has similar checks for benchmarks and TB interface.
- [tparallel](https://github.com/moricho/tparallel) finds inappropriate usage of `t.Parallel()` method in your Go test codes.
- [unparam](https://github.com/mvdan/unparam) reports unused function parameters and results in your code.
- [unused](https://github.com/dominikh/go-tools/tree/master/unused) finds unused code.
- [usestdlibvars](https://github.com/sashamelentyev/usestdlibvars) detects the possibility to use variables/constants from the Go standard library.
- [varnamelen](https://github.com/blizzy78/varnamelen) checks that the length of a variable's name matches its usage scope.
- [wastedassign](https://github.com/sanposhiho/wastedassign) finds wasted assignment statements.

## Some other linters

- [go vet](https://pkg.go.dev/cmd/vet)
- [golangci-linter](https://golangci-lint.run)
- [revive](https://revive.run)

## License

MIT licensed. See the bundled [LICENSE](LICENSE) file for more details.
