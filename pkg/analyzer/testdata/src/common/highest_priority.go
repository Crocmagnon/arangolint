package common

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

type svc struct {
	opts *arangodb.BeginTransactionOptions
}

type nested struct {
	conf struct {
		txnOpts *arangodb.BeginTransactionOptions
	}
}

func highestPriorityCases() {
	ctx := context.Background()
	arangoClient := arangodb.NewClient(nil)
	db, _ := arangoClient.GetDatabase(ctx, "name", nil)

	// SelectorExpr as argument
	s := &svc{}
	s.opts = &arangodb.BeginTransactionOptions{}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, s.opts) // want "missing AllowImplicit option"
	s.opts.AllowImplicit = true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, s.opts) // ok

	// Nested selector
	n := &nested{}
	n.conf.txnOpts = &arangodb.BeginTransactionOptions{}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, n.conf.txnOpts) // want "missing AllowImplicit option"
	n.conf.txnOpts.AllowImplicit = true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, n.conf.txnOpts) // ok

	// ParenExpr around identifier
	var p *arangodb.BeginTransactionOptions
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, (p)) // want "missing AllowImplicit option"

	// Typed conversion to pointer type with nil
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, (*arangodb.BeginTransactionOptions)(nil)) // want "missing AllowImplicit option"
}
