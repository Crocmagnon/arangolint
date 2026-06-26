# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`arangolint` is a static analyzer (linter) for ArangoDB's Go driver v2 that provides two main features:

1. **Transaction Safety**: Enforces the explicit use of the `AllowImplicit` option when calling `BeginTransaction()` to prevent deadlocks from lazily-added transaction collections.
2. **AQL Injection Prevention**: Detects potentially vulnerable AQL query construction using string concatenation or `fmt.Sprintf`, encouraging the use of bind variables instead.

The analyzer is integrated into `golangci-lint` since v2.2.0.

## Architecture

### Core Components

**Main Analyzer (`pkg/analyzer/analyzer.go`)**
The analyzer is built on Go's `analysis` framework from `golang.org/x/tools/go/analysis`. Key design decisions:

- **Intra-procedural only**: Does not follow values across function/method boundaries
- **Flow- and block-sensitive**: Scans statements in lexical order within the current function, considering only prior statements in the nearest block and ancestor blocks
- **Conservative by design**: Assumes `AllowImplicit` is set when options come from unknown factory/helper calls to avoid false positives

**Detection Logic**

The analyzer provides two types of detection:

**1. Transaction AllowImplicit Detection**

Identifies calls to `arangodb.Database.BeginTransaction()` and checks the third parameter (options):

1. Direct composite literals: `&arangodb.BeginTransactionOptions{AllowImplicit: true}`
2. Variable assignments before the call: `opts.AllowImplicit = true`
3. Initialization in variable declarations
4. Control-flow structures (if/else, for, range, switch)
5. Package-level variable initializations

**2. AQL Query Injection Detection**

Identifies calls to query methods (`Query`, `QueryBatch`, `ValidateQuery`, `ExplainQuery`) on both `Database` and `Transaction` types and analyzes the query string parameter:

1. Direct string concatenation: `"FOR u IN users FILTER u.name == '" + userName + "' RETURN u"`
2. `fmt.Sprintf` calls: `fmt.Sprintf("FOR u IN users FILTER u.name == '%s' RETURN u", userName)`
3. Variable assignments with concatenation: `query := "FOR u IN users" + userFilter`
4. Compound assignment operators: `query += " FILTER u.name == '" + userName + "'"`
5. Control-flow structures (if/else, for, range, switch) that build queries with concatenation
6. Package-level variable initializations with concatenation

The analyzer distinguishes between unsafe concatenation (involving variables) and safe static concatenation (only string literals).

**Type Resolution**
Uses `types.AssignableTo()` to detect calls through embedded types, wrappers, or aliases. Falls back to string suffix matching for type names when full type information is unavailable.

**Index Expression Handling**
Tracks assignments to array/slice elements (e.g., `arr[i].AllowImplicit = true`) by matching both the base identifier and constant index values.

### CLI Entrypoint

`cmd/arangolint/main.go` uses `singlechecker.Main()` to provide a standalone analyzer binary that integrates with standard Go tooling.

### Test Structure

Tests use `analysistest.Run()` with golden test files:
- `pkg/analyzer/testdata/src/common/` — Standard tests covering various patterns, one file per pattern family. Examples:
  - Query injection detection (`query_injection.go`, `transaction_query_injection.go`)
  - Transaction `AllowImplicit` detection across control flow, shadowing, embeddings/wrappers, type aliases, etc. (`shadowing.go`, `embeddings_and_wrappers.go`, `for_loops.go`, `switch_cases.go`, …)
- `pkg/analyzer/testdata/src/cgo/` — CGO-specific test cases (`cgo.go`)
- Each test file includes `// want "..."` comments marking expected diagnostics
- Test data includes vendored copies of ArangoDB driver v2 for self-contained testing

## Common Commands

### Build
```bash
go build ./...
```

### Run Tests
```bash
mise run test   # runs: go test -race -v ./...
# or directly (CGO_ENABLED=1 is required for the cgo test cases):
CGO_ENABLED=1 go test ./...
```

### Run Linter
```bash
mise run lint   # runs: golangci-lint run --fix
# or directly:
CGO_ENABLED=1 golangci-lint run ./...
```

### Install golangci-lint
```bash
mise install
```
`mise.toml` pins the `golangci-lint` version; `mise install` provisions it.

### Run Single Test
```bash
CGO_ENABLED=1 go test ./pkg/analyzer -run TestAnalyzer/common
CGO_ENABLED=1 go test ./pkg/analyzer -run TestAnalyzer/cgo
```

### Tidy Dependencies (including test data)
```bash
mise run tidy
```
This tidies the main module and the two test data modules (`cgo` and `common`), re-vendoring their dependencies.

## Module & Toolchain

- **Module path**: `go.augendre.info/arangolint` (does not match the on-disk directory path)
- **Go**: `1.25.0`; the only notable direct dependency is `golang.org/x/tools` (provides the `analysis` framework)
- **`mise.toml`** pins `golangci-lint 2.12.2`, `goreleaser 2`, and `prek` (pre-commit runner), and sets `CGO_ENABLED=1` for the dev environment. It also defines a `release-snapshot` task (`goreleaser release --snapshot --clean`); releases are configured in `.goreleaser.yaml`
- **`.pre-commit-config.yaml`** runs golangci-lint plus `go test` on commit
- **`.golangci.yml`** (config v2, `fix: true`) enables nearly all linters except `depguard`, `exhaustruct`, `nonamedreturns`, and `wsl`; `cyclop` caps cyclomatic complexity at 15. Formatting is done by `goimports`/`gofmt`/`gofumpt`/`golines` with local prefix `go.augendre.info/arangolint`. Match this style to pass lint on the first try.

## Development Notes

### Adding New Detection Patterns

When extending the analyzer to handle new patterns:

1. Add test cases to `pkg/analyzer/testdata/src/common/` or `pkg/analyzer/testdata/src/cgo/`
2. Use `// want "..."` comments to mark expected diagnostics (e.g., `// want "missing AllowImplicit option"` or `// want "query string uses concatenation"`)
3. Implement detection logic in the appropriate function:
   - For `AllowImplicit` detection: `shouldReportMissingAllowImplicit()` or related helpers
   - For query injection detection: `shouldReportQueryConcatenation()` or related helpers
4. Maintain conservative behavior: when uncertain, do not report to avoid false positives
5. Consider both direct patterns (composite literals, direct concatenation) and indirect patterns (variable assignments, control-flow structures)

### Analyzing Call Sites

The analyzer uses `inspector.WithStack()` to traverse call expressions with their enclosing blocks. The stack provides context for flow-sensitive analysis through `ancestorBlocks()` and `scanPriorStatements()`.

### Root Identifier Resolution

`rootIdent()` peels nested expressions (parens, stars, selectors, index/slice) to find the underlying identifier. This enables tracking values through dereferences, field accesses, and array indexing.

### Query Injection Detection Helpers

Key functions for query injection analysis:

- `identifyQueryMethod()`: Identifies calls to query methods (`Query`, `QueryBatch`, `ValidateQuery`, `ExplainQuery`) on both `Database` and `Transaction` types and returns the query argument index
- `getQueryArgIndex()`: Maps method names to their query argument index
- `isQueryReceiverType()`: Checks if a type is `Database` or `Transaction` using type resolution
- `getArangoDBTypes()`: Retrieves `Database` and `Transaction` types from the arangodb package
- `lookupType()`: Helper to lookup types by name in a package scope
- `isConcatenatedString()`: Recursively checks if a binary expression uses `+` with at least one non-literal operand
- `isAllStringLiterals()`: Distinguishes safe static concatenation from unsafe variable interpolation
- `isFmtSprintfCall()`: Detects `fmt.Sprintf` calls by checking package origin
- `wasBuiltWithConcatenation()`: Traces back through assignments to determine if a variable was built with concatenation
- `stmtAssignsConcatenation()`: Checks if a statement assigns a concatenated query to a given variable, including support for control-flow structures

The concatenation detection is flow-sensitive and traces through variable assignments, declarations, and control structures, but remains intra-procedural for performance and to avoid false positives.

### Conservative Analysis Philosophy

The analyzer prioritizes avoiding false positives over catching every violation. Unknown patterns are assumed valid:
- For `AllowImplicit`: Options from helper functions are assumed to have `AllowImplicit` set
- For query injection: Queries from helper functions are not flagged since they may use bind variables internally

This keeps the linter practical for real-world codebases where values may be constructed indirectly through abstractions.
