package common

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

func usePackageVars() {
	ctx := context.Background()
	arangoClient := arangodb.NewClient(nil)
	db, _ := arangoClient.GetDatabase(ctx, "name", nil)

	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, UnsafePkgVar) // want "missing AllowImplicit option"
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, SafePkgVar)
}
