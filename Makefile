.PHONY: tidy install-linter test lint
tidy:
	go mod tidy
	cd pkg/analyzer/testdata/src/cgo && go mod tidy && go mod vendor
	cd pkg/analyzer/testdata/src/common && go mod tidy && go mod vendor
test:
	CGO_ENABLED=1 go test ./...
lint:
	CGO_ENABLED=1 golangci-lint run ./...
install-linter:
	@./install-linter
