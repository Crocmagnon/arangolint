package common

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

func rangeLoops() {
	ctx := context.Background()
	arangoClient := arangodb.NewClient(nil)
	db, _ := arangoClient.GetDatabase(ctx, "name", nil)

	// Mutate inside a range loop before a call after the loop.
	opts := &arangodb.BeginTransactionOptions{}
	for range []int{1, 2, 3} {
		opts.AllowImplicit = true
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts) // want "missing AllowImplicit option"

	// Control: no mutation in range.
	opts2 := &arangodb.BeginTransactionOptions{}
	for range []int{} {
		// nothing
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts2) // want "missing AllowImplicit option"
}
