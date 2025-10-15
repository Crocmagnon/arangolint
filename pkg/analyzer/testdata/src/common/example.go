package common

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

type dbclient struct {
	db arangodb.Database
}

func example() {
	ctx := context.Background()
	arangoClient := arangodb.NewClient(nil)
	db, _ := arangoClient.GetDatabase(ctx, "name", nil)
	dbc := &dbclient{db: db}

	// direct nil
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, nil)           // want "missing AllowImplicit option"
	dbc.db.BeginTransaction(ctx, arangodb.TransactionCollections{}, nil)       // want "missing AllowImplicit option"
	trx, _ := db.BeginTransaction(ctx, arangodb.TransactionCollections{}, nil) // want "missing AllowImplicit option"
	_ = trx

	// direct missing
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{LockTimeout: 0})          // want "missing AllowImplicit option"
	dbc.db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{LockTimeout: 0})      // want "missing AllowImplicit option"
	trx, _ = db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{LockTimeout: 0}) // want "missing AllowImplicit option"

	// direct false
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: false})
	dbc.db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: false})
	trx, _ = db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: false})

	// direct true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: true})
	dbc.db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: true})
	trx, _ = db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: true})

	// direct with other fields
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: true, LockTimeout: 0})
	dbc.db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: true, LockTimeout: 0})
	trx, _ = db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: true, LockTimeout: 0})

	// indirect no pointer
	options := arangodb.BeginTransactionOptions{LockTimeout: 0}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &options) // want "missing AllowImplicit option"
	options.AllowImplicit = true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &options)

	// indirect pointer
	optns := &arangodb.BeginTransactionOptions{LockTimeout: 0}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, optns) // want "missing AllowImplicit option"
	optns.AllowImplicit = true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, optns)

	// var declaration (no pointer)
	var options2 = arangodb.BeginTransactionOptions{LockTimeout: 0}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &options2) // want "missing AllowImplicit option"
	options2.AllowImplicit = true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &options2)

	// var declaration (pointer)
	var optns2 = &arangodb.BeginTransactionOptions{LockTimeout: 0}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, optns2) // want "missing AllowImplicit option"
	optns2.AllowImplicit = true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, optns2)

	// var declaration with AllowImplicit in init (no pointer)
	var options3 = arangodb.BeginTransactionOptions{AllowImplicit: true}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &options3)

	// var declaration with AllowImplicit in init (pointer)
	var optns3 = &arangodb.BeginTransactionOptions{AllowImplicit: true}
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, optns3)
}
