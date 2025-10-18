.PHONY: tidy install-linter test lint
tidy:
	go mod tidy
	cd pkg/analyzer/testdata/src/cgo && go mod tidy && go mod vendor
	cd pkg/analyzer/testdata/src/common && go mod tidy && go mod vendor
test:
	go test ./...
lint:
	golangci-lint run ./...
install-linter:
	@./install-linter
