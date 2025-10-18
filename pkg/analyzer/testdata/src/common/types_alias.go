package common

import (
	"github.com/arangodb/go-driver/v2/arangodb"
)

type T = arangodb.Database
type T2 arangodb.Database
type T3 interface {
	arangodb.Database
}

func f2() {
	var t T
	t.BeginTransaction(nil, arangodb.TransactionCollections{}, nil) // want "missing AllowImplicit option"
	var t2 T2
	t2.BeginTransaction(nil, arangodb.TransactionCollections{}, nil) // want "missing AllowImplicit option"
	var t3 T3
	t3.BeginTransaction(nil, arangodb.TransactionCollections{}, nil) // want "missing AllowImplicit option"
}
