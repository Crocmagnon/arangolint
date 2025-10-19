package common

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

func switchCases() {
	ctx := context.Background()
	arangoClient := arangodb.NewClient(nil)
	db, _ := arangoClient.GetDatabase(ctx, "name", nil)

	// 1) Assign in switch init clause; call after switch.
	opts := &arangodb.BeginTransactionOptions{}
	switch opts = (&arangodb.BeginTransactionOptions{AllowImplicit: true}); 0 {
	default:
		// no-op
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts)

	// 2) Assign in a specific case; call after switch.
	opts2 := &arangodb.BeginTransactionOptions{}
	switch 1 {
	case 1:
		opts2.AllowImplicit = true
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts2)

	// 3) Default-only assignment; call after switch.
	opts3 := &arangodb.BeginTransactionOptions{}
	switch 2 {
	default:
		opts3.AllowImplicit = true
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts3)

	// 4) Control: switch without any assignment; call after switch should be flagged.
	opts4 := &arangodb.BeginTransactionOptions{}
	switch 3 {
	case 1:
		// no assignment
	case 2:
		// no assignment
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts4) // want "missing AllowImplicit option"

	// 5) Fallthrough between cases with and without assignment; conservative behavior should treat any case assignment before call as sufficient.
	opts5 := &arangodb.BeginTransactionOptions{}
	switch 4 {
	case 1:
		// no assignment
		fallthrough
	case 2:
		opts5.AllowImplicit = true
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts5)
}
