package common

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

type wrapper struct {
	arangodb.Database
}

func embeddingsAndWrappers() {
	ctx := context.Background()
	arangoClient := arangodb.NewClient(nil)
	db, _ := arangoClient.GetDatabase(ctx, "name", nil)

	var w wrapper
	w.Database = db

	// Using embedded arangodb.Database; analyzer should detect BeginTransaction on wrapper.
	w.BeginTransaction(ctx, arangodb.TransactionCollections{}, nil) // want "missing AllowImplicit option"

	opts := &arangodb.BeginTransactionOptions{AllowImplicit: true}
	w.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts)
}
