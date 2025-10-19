package common

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

func forLoops() {
	ctx := context.Background()
	arangoClient := arangodb.NewClient(nil)
	db, _ := arangoClient.GetDatabase(ctx, "name", nil)

	// 1) Assign in init clause before the loop; call occurs after the loop.
	// Use assignment (not short var) so opts is in the outer scope.
	var i int
	opts := &arangodb.BeginTransactionOptions{}
	for opts = (&arangodb.BeginTransactionOptions{AllowImplicit: true}); i < 1; i++ {
		// no-op
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts)

	// 2) Assign inside loop body; call occurs after the loop.
	opts2 := &arangodb.BeginTransactionOptions{}
	for j := 0; j < 1; j++ {
		opts2.AllowImplicit = true
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts2)

	// 3) Control: same shape without any assignment in init/body should be flagged when calling after the loop.
	opts3 := &arangodb.BeginTransactionOptions{}
	for k := 0; k < 1; k++ {
		// no assignment to AllowImplicit
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts3) // want "missing AllowImplicit option"

	// 4) Nested for loops: inner loop sets, call after outer loop.
	opts4 := &arangodb.BeginTransactionOptions{}
	for k := 0; k < 1; k++ {
		for m := 0; m < 1; m++ {
			opts4.AllowImplicit = true
		}
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts4)
}
