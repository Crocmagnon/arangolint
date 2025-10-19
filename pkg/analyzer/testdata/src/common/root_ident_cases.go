package common

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

// Reuse svc and nested types from highest_priority.go
var globalSvc = &svc{}

func getGlobalSvc() *svc { return globalSvc }

func rootIdentCases() {
	ctx := context.Background()
	arangoClient := arangodb.NewClient(nil)
	db, _ := arangoClient.GetDatabase(ctx, "name", nil)

	// 1) Selector with parenthesized root: (s).opts
	s := &svc{}
	s.opts = &arangodb.BeginTransactionOptions{}
	s.opts.AllowImplicit = true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, (s).opts)

	// 2) Selector with dereferenced root: (*s2).opts
	s2 := &svc{}
	s2.opts = &arangodb.BeginTransactionOptions{}
	(*s2).opts.AllowImplicit = true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, (*s2).opts)

	// 3) Nested selector with parens and star: (*(n)).conf.txnOpts
	n := &nested{}
	n.conf.txnOpts = &arangodb.BeginTransactionOptions{}
	(*(n)).conf.txnOpts.AllowImplicit = true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, (*(n)).conf.txnOpts)

	// 4) Non-ident root via CallExpr: getGlobalSvc().opts — even if assigned earlier, should be flagged
	getGlobalSvc().opts = &arangodb.BeginTransactionOptions{}
	getGlobalSvc().opts.AllowImplicit = true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, getGlobalSvc().opts) // want "missing AllowImplicit option"

	// 5) Non-ident root via IndexExpr: arr[0].opts — even if assigned earlier, should be flagged
	arr := make([]svc, 1)
	arr[0].opts = &arangodb.BeginTransactionOptions{}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, arr[0].opts) // want "missing AllowImplicit option"
	arr[0].opts.AllowImplicit = true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, arr[0].opts)

	arr2 := make([]*arangodb.BeginTransactionOptions, 1)
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, arr2[0]) // want "missing AllowImplicit option"
	arr2[0].AllowImplicit = true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, arr2[0])

	arr3 := make([]arangodb.BeginTransactionOptions, 1)
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arr3[0]) // want "missing AllowImplicit option"
	arr3[0].AllowImplicit = true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arr3[0])

	arr4 := make([]arangodb.BeginTransactionOptions, 2)
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arr4[0]) // want "missing AllowImplicit option"
	arr4[1].AllowImplicit = true                                          // updating 1 not 0
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arr4[0]) // want "missing AllowImplicit option"
}
