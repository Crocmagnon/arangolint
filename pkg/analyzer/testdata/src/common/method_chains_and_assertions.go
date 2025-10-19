package common

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

func getDB(ctx context.Context) arangodb.Database {
	client := arangodb.NewClient(nil)
	db, _ := client.GetDatabase(ctx, "name", nil)
	return db
}

func methodChainsAndAssertions() {
	ctx := context.Background()
	db := getDB(ctx)

	// Basic call via helper expression result
	opts := &arangodb.BeginTransactionOptions{}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts) // want "missing AllowImplicit option"
	opts.AllowImplicit = true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts)

	// Type assertion to arangodb.Database before calling BeginTransaction
	var dbi any = db
	opts2 := &arangodb.BeginTransactionOptions{}
	dbi.(arangodb.Database).BeginTransaction(ctx, arangodb.TransactionCollections{}, opts2) // want "missing AllowImplicit option"
	opts2.AllowImplicit = true
	dbi.(arangodb.Database).BeginTransaction(ctx, arangodb.TransactionCollections{}, opts2)
}
