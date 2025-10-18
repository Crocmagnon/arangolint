package common

import "context"

type SomeOtherDB struct{}

func (s SomeOtherDB) BeginTransaction(_ context.Context, _, _ any) {}

func f() {
	var db SomeOtherDB
	db.BeginTransaction(nil, nil, nil)
}
