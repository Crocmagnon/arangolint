# arangolint
[![Go Reference](https://pkg.go.dev/badge/github.com/Crocmagnon/arangolint.svg)](https://pkg.go.dev/github.com/Crocmagnon/arangolint)
[![Go Report Card](https://goreportcard.com/badge/github.com/Crocmagnon/arangolint)](https://goreportcard.com/report/github.com/Crocmagnon/arangolint)
[![Go Coverage](https://github.com/Crocmagnon/arangolint/wiki/coverage.svg)](https://github.com/Crocmagnon/arangolint/wiki/Coverage)

Opinionated linter for [ArangoDB go driver v2](https://github.com/arangodb/go-driver).

`arangolint` is available in `golangci-lint` since v2.2.0.

## Features

### Enforce explicit `AllowImplicit` in transactions
Why? Because it forces you as a developer to evaluate the need of implicit collections in transactions.

Why should you? Because [lazily adding collections](https://docs.arangodb.com/3.11/develop/transactions/locking-and-isolation/#lazily-adding-collections) to transactions can lead to deadlocks, and because the default is to allow it.

```go
ctx := context.Background()
arangoClient := arangodb.NewClient(nil)
db, _ := arangoClient.GetDatabase(ctx, "name", nil)

// Bad
trx, _ := db.BeginTransaction(ctx, arangodb.TransactionCollections{}, nil) // want "missing AllowImplicit option"
trx, _ = db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{LockTimeout: 0}) // want "missing AllowImplicit option"

// Good
trx, _ = db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: true})
trx, _ = db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: false})
trx, _ = db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: true, LockTimeout: 0})

// Indirect via variable (no pointer)
options := arangodb.BeginTransactionOptions{LockTimeout: 0}
db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &options) // want "missing AllowImplicit option"
options.AllowImplicit = true
db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &options)

// Indirect via pointer variable
optns := &arangodb.BeginTransactionOptions{LockTimeout: 0}
db.BeginTransaction(ctx, arangodb.TransactionCollections{}, optns) // want "missing AllowImplicit option"
optns.AllowImplicit = true
db.BeginTransaction(ctx, arangodb.TransactionCollections{}, optns)
```

Notes and limitations:
* Variable tracking is block-scoped and flow-sensitive across the nearest and ancestor blocks within the current function.
* It detects AllowImplicit when set in the composite literal initialization or via an explicit assignment (e.g., options.AllowImplicit = ...).
* It does not perform inter-procedural analysis or track values across complex control flow at this time.
