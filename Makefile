all: tidy fix fmt generate-config

fmt:
	@go fmt ./...

fix:
	@go run ./... -config gochecker.yaml -fix ./...

vet:
	@go run ./... -config gochecker.yaml ./...

github-vet:
	@go run ./... -config gochecker.yaml -output github ./...

test-all:
	@go run ./... -config config_examples/config.yaml ./...

generate-config:
	@go run ./... generate-config > config_examples/config.yaml

tidy:
	@go mod tidy

godoc:
	@go install golang.org/x/tools/cmd/godoc@latest
	@echo http://localhost:6060/pkg/github.com/sv-tools/openapi/
	@godoc -http=:6060 >/dev/null
