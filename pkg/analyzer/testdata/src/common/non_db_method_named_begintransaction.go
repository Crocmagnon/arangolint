package common

import "context"

type otherDBPtr struct{}

func (s *otherDBPtr) BeginTransaction(_ context.Context, _, _ any) {}

// Another variant with wrong arity to ensure analyzer ignores non-matching signatures.
type oddArityDB struct{}

func (oddArityDB) BeginTransaction(_, _ any) {}

func nonDBMethods() {
	var db1 otherDBPtr
	db1.BeginTransaction(nil, nil, nil)

	var db2 oddArityDB
	db2.BeginTransaction(nil, nil)
}
