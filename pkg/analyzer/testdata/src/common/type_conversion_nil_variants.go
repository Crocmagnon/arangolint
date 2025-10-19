package common

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

// optsFactory returns a *arangodb.BeginTransactionOptions, used to exercise the
// isTypeConversionToTxnOptionsPtrNil fallback path where the 3rd argument is a
// regular call expression rather than a type conversion. The analyzer should be
// conservative and not flag this shape.
func optsFactory(_ any) *arangodb.BeginTransactionOptions { return nil }

func typeConversionNilVariants() {
	ctx := context.Background()
	client := arangodb.NewClient(nil)
	db, _ := client.GetDatabase(ctx, "name", nil)

	// 1) Positive: deep-parenthesized pointer-type conversion to nil should be flagged.
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, (((*arangodb.BeginTransactionOptions))(nil))) // want "missing AllowImplicit option"

	// 2) Negative: third arg is a regular function call with a single nil argument,
	// returning *arangodb.BeginTransactionOptions. Analyzer should not flag.
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, optsFactory(nil))


	// 4) Negative: pointer-type conversion with parenthesized nil; current behavior
	// does not treat (nil) as a bare nil ident. Should not flag.
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, (*arangodb.BeginTransactionOptions)((nil)))
}
