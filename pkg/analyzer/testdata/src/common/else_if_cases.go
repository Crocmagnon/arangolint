package common

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

func elseIfCases() {
	ctx := context.Background()
	arangoClient := arangodb.NewClient(nil)
	db, _ := arangoClient.GetDatabase(ctx, "name", nil)

	// 1) Else branch sets; call after the if.
	opts := &arangodb.BeginTransactionOptions{}
	if false {
		// no assignment
	} else {
		opts.AllowImplicit = true
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts)

	// 2) Else-if branch sets; call after the chain.
	opts2 := &arangodb.BeginTransactionOptions{}
	if false {
		// no
	} else if true {
		opts2.AllowImplicit = true
	} else {
		// no
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts2)

	// 3) If branch sets; call after the chain.
	opts3 := &arangodb.BeginTransactionOptions{}
	if true {
		opts3.AllowImplicit = true
	} else {
		// no
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts3)

	// 4) Control: none of the branches set; expect diagnostic.
	opts4 := &arangodb.BeginTransactionOptions{}
	if false {
		// no
	} else if false {
		// no
	} else {
		// no
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts4) // want "missing AllowImplicit option"

	// 5) Nested else-if where only final else sets; call after.
	opts5 := &arangodb.BeginTransactionOptions{}
	if false {
		// no
	} else if false {
		// no
	} else {
		opts5.AllowImplicit = true
	}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, opts5)
}
