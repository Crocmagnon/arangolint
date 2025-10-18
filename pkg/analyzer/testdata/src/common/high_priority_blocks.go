package common

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

func highPriorityBlocks() {
	ctx := context.Background()
	arangoClient := arangodb.NewClient(nil)
	db, _ := arangoClient.GetDatabase(ctx, "name", nil)

	// 1) Mutation in earlier nested block before call
	opts := &arangodb.BeginTransactionOptions{}
	if true {
		opts.AllowImplicit = true
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts)

	// Control: without mutation should be flagged
	opts2 := &arangodb.BeginTransactionOptions{}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts2) // want "missing AllowImplicit option"

	// 3) Mutation in init clause before later call outside the if
	var opts4 = &arangodb.BeginTransactionOptions{}
	if x := 1; x > 0 {
		opts4.AllowImplicit = true
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts4)
}
