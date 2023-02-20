# gochecker
Go tool to check the code (linter)

Another variation of the go linters.
The `gochecker` is a wrapper for [multichecker](https://pkg.go.dev/golang.org/x/tools/go/analysis/multichecker).
The `gochecker` supports `go vet` interface and includes all official [analiyzers](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes)
and some custom linters, see [analyzers.go](analyzers.go) file for the full list.
In theory the `gochecker` can support any analyzers which implement [Analyzer type](https://pkg.go.dev/golang.org/x/tools/go/analysis#Analyzer) 

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
gochecker -fieldalignment ./...
```

and please check `gochecker help` or `gochecker help <analyzer>` for full help.

## TODO

- [x] go analyzers
- [ ] custom analyzers (linters)
- [ ] `diff` output
- [ ] `github` output
- [ ] `nolint` and/or global exclude rules

## Some other linters

- [go vet](https://pkg.go.dev/cmd/vet)
- [golangci-linter](https://golangci-lint.run)
- [revive](https://revive.run)

## License

MIT licensed. See the bundled [LICENSE](LICENSE) file for more details.
