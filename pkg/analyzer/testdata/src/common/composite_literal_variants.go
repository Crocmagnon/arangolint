package common

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

func compositeLiteralVariants() {
	ctx := context.Background()
	arangoClient := arangodb.NewClient(nil)
	db, _ := arangoClient.GetDatabase(ctx, "name", nil)

	// Address-of composite with AllowImplicit not as first element and spread across multiple lines with comments
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{
		LockTimeout: 1,
		// important flag below
		AllowImplicit: true,
	})

	// Mixed order of other fields before/after AllowImplicit, including trailing comma
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{
		AllowImplicit: true,
		LockTimeout:  0,
	})
}
