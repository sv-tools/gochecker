all: tidy fmt generate-config lint

fmt:
	@go fmt ./...

lint:
	@go run ./... -config config_examples/vet_config.yaml ./...

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