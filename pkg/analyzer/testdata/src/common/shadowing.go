package common

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

func shadowingAndOrder() {
	ctx := context.Background()
	arangoClient := arangodb.NewClient(nil)
	db, _ := arangoClient.GetDatabase(ctx, "name", nil)

	// 1) Inner-scope shadowing of options variable: expect diagnostic for inner variable.
	opts := &arangodb.BeginTransactionOptions{}
	opts.AllowImplicit = true
	{
		opts := &arangodb.BeginTransactionOptions{}
		db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts) // want "missing AllowImplicit option"
	}

	// 2) Mutation after call should not count.
	opts2 := &arangodb.BeginTransactionOptions{}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts2) // want "missing AllowImplicit option"
	opts2.AllowImplicit = true

	// 3) Multiple assignments before call (true then false): presence of field in any assignment is enough.
	opts3 := &arangodb.BeginTransactionOptions{}
	opts3.AllowImplicit = true
	opts3.AllowImplicit = false
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts3)
}
