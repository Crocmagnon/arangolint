package common

import "github.com/arangodb/go-driver/v2/arangodb"

// Package-level variables declared in a separate file
var UnsafePkgVar = &arangodb.BeginTransactionOptions{}
var SafePkgVar = &arangodb.BeginTransactionOptions{AllowImplicit: true}
