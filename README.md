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

// Indirect
options := arangodb.BeginTransactionOptions{LockTimeout: 0}
db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &options) // want "missing AllowImplicit option"
options.AllowImplicit = true
db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &options)
```

Notes and limitations:
- Intra-procedural only: the analyzer does not follow values across function/method boundaries.
- Flow- and block-sensitive within the current function: prior statements in the nearest block and its ancestor blocks are considered when evaluating a call site.
- What is detected: AllowImplicit set either in a composite literal initialization (e.g., &arangodb.BeginTransactionOptions{AllowImplicit: true}) or via an explicit assignment before the call (e.g., opts.AllowImplicit = ...).
- Conservative by design: when the options value comes from an unknown factory/helper call, arangolint assumes AllowImplicit may be set to avoid false positives.
- Out of scope (for now): inter-procedural tracking, deep control-flow analysis, and inference through complex aliasing beyond simple identifiers and selectors.

### Detect AQL query injection vulnerabilities

Why? Because building AQL queries using string concatenation or `fmt.Sprintf` with user-supplied values can lead to AQL injection attacks, similar to SQL injection vulnerabilities.

Why should you use bind variables? Because [bind variables](https://docs.arangodb.com/3.11/aql/how-to-invoke-aql/with-arangosh/#bind-parameters) are automatically escaped and sanitized by ArangoDB, preventing injection attacks.

```go
ctx := context.Background()
db, _ := arangoClient.GetDatabase(ctx, "name", nil)
userName := "admin"
userAge := 25

// Bad - String concatenation with variables
db.Query(ctx, "FOR u IN users FILTER u.name == '"+userName+"' RETURN u", nil) // want "query string uses concatenation"

// Bad - Building query with multiple concatenations
query := "FOR u IN users FILTER u.name == '" + userName + "' RETURN u"
db.Query(ctx, query, nil) // want "query string uses concatenation"

// Bad - Using fmt.Sprintf
sprintfQuery := fmt.Sprintf("FOR u IN users FILTER u.name == '%s' RETURN u", userName)
db.Query(ctx, sprintfQuery, nil) // want "query string uses concatenation"

// Bad - Compound assignment +=
q := "FOR u IN users"
q += " FILTER u.name == '" + userName + "'"
db.Query(ctx, q, nil) // want "query string uses concatenation"

// Good - Using bind variables
db.Query(ctx, "FOR u IN users FILTER u.name == @name RETURN u", &arangodb.QueryOptions{
    BindVars: map[string]interface{}{
        "name": userName,
    },
})

// Good - Multiple bind variables
db.Query(ctx, "FOR u IN users FILTER u.name == @name AND u.age == @age RETURN u", &arangodb.QueryOptions{
    BindVars: map[string]interface{}{
        "name": userName,
        "age":  userAge,
    },
})

// Good - Static query string (no variable interpolation)
db.Query(ctx, "FOR u IN users RETURN u", nil)

// Good - Static concatenation (only string literals)
staticQuery := "FOR u IN users" + " FILTER u.age > 18" + " RETURN u"
db.Query(ctx, staticQuery, nil)

// Works with transactions too
trx, _ := db.BeginTransaction(ctx, arangodb.TransactionCollections{
    Read: []string{"users"},
}, &arangodb.BeginTransactionOptions{AllowImplicit: false})

// Bad - Transaction query with concatenation
trx.Query(ctx, "FOR u IN users FILTER u.name == '"+userName+"' RETURN u", nil) // want "query string uses concatenation"

// Good - Transaction query with bind variables
trx.Query(ctx, "FOR u IN users FILTER u.name == @name RETURN u", &arangodb.QueryOptions{
    BindVars: map[string]interface{}{
        "name": userName,
    },
})
```

Covered methods on `Database` and `Transaction`:
- `Query()`
- `QueryBatch()`
- `ValidateQuery()`
- `ExplainQuery()`

Notes and limitations:
- Intra-procedural only: the analyzer does not follow values across function/method boundaries.
- Detects direct concatenation (`+` operator) and `fmt.Sprintf` calls in the same function.
- Detects concatenation in variable assignments, declarations, and control-flow structures (if/else, for, range, switch).
- Conservative by design: queries from helper functions or factory methods are not flagged to avoid false positives.
- Static string concatenation (only literals, no variables) is considered safe and not flagged.
