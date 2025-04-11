# arangolint
[![Go Reference](https://pkg.go.dev/badge/github.com/Crocmagnon/arangolint.svg)](https://pkg.go.dev/github.com/Crocmagnon/arangolint)
[![Go Report Card](https://goreportcard.com/badge/github.com/Crocmagnon/arangolint)](https://goreportcard.com/report/github.com/Crocmagnon/arangolint)
[![Go Coverage](https://github.com/Crocmagnon/arangolint/wiki/coverage.svg)](https://github.com/Crocmagnon/arangolint/wiki/Coverage)

Opinionated linter for ArangoDB go driver

## Features

### Enforce explicit `AllowImplicit` in transactions
Why? Because it forces you as a developer to evaluate the need of implicit collections in transactions.

Why should you? Because [lazily adding collections](https://docs.arangodb.com/3.11/develop/transactions/locking-and-isolation/#lazily-adding-collections) to transactions can lead to deadlocks, and because the default if this value is not specified is to allow it.

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
```

Limitations: this currently only works when transaction options are directly passed to `BeginTransaction`, not when using a variable.
