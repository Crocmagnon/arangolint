package common

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver/v2/arangodb"
)

func transactionQueryInjectionExamples(db arangodb.Database, userName string, userAge int) {
	ctx := context.Background()

	// Create a transaction
	trx, _ := db.BeginTransaction(ctx, arangodb.TransactionCollections{
		Read:  []string{"users"},
		Write: []string{"users"},
	}, &arangodb.BeginTransactionOptions{AllowImplicit: false})

	// UNSAFE: Direct string concatenation with +
	trx.Query(ctx, "FOR u IN users FILTER u.name == '"+userName+"' RETURN u", nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: String concatenation across multiple lines
	query := "FOR u IN users FILTER u.name == '" + userName + "' RETURN u"
	trx.Query(ctx, query, nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: Multiple concatenations
	trx.Query(ctx, "FOR u IN users FILTER u.name == '"+userName+"' AND u.age == "+fmt.Sprint(userAge)+" RETURN u", nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: Concatenation with compound assignment +=
	q := "FOR u IN users"
	q += " FILTER u.name == '" + userName + "'"
	q += " RETURN u"
	trx.Query(ctx, q, nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: fmt.Sprintf
	sprintfQuery := fmt.Sprintf("FOR u IN users FILTER u.name == '%s' RETURN u", userName)
	trx.Query(ctx, sprintfQuery, nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: QueryBatch with concatenation
	trx.QueryBatch(ctx, "FOR u IN users FILTER u.name == '"+userName+"' RETURN u", nil, nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: ValidateQuery with concatenation
	trx.ValidateQuery(ctx, "FOR u IN users FILTER u.name == '"+userName+"' RETURN u") // want "query string uses concatenation instead of bind variables"

	// UNSAFE: ExplainQuery with concatenation
	trx.ExplainQuery(ctx, "FOR u IN users FILTER u.name == '"+userName+"' RETURN u", nil, nil) // want "query string uses concatenation instead of bind variables"

	// SAFE: Using bind variables properly
	trx.Query(ctx, "FOR u IN users FILTER u.name == @name RETURN u", &arangodb.QueryOptions{
		BindVars: map[string]interface{}{
			"name": userName,
		},
	})

	// SAFE: Using bind variables with multiple parameters
	trx.Query(ctx, "FOR u IN users FILTER u.name == @name AND u.age == @age RETURN u", &arangodb.QueryOptions{
		BindVars: map[string]interface{}{
			"name": userName,
			"age":  userAge,
		},
	})

	// SAFE: Static query string (no variables)
	trx.Query(ctx, "FOR u IN users RETURN u", nil)

	// SAFE: Static query with options but no variables
	trx.Query(ctx, "FOR u IN users FILTER u.age > 18 RETURN u", &arangodb.QueryOptions{
		Count: true,
	})

	// SAFE: Concatenation of only static strings (no variables) - this is OK
	staticQuery := "FOR u IN users" + " FILTER u.age > 18" + " RETURN u"
	trx.Query(ctx, staticQuery, nil)

	// SAFE: QueryBatch with bind vars
	trx.QueryBatch(ctx, "FOR u IN users FILTER u.name == @name RETURN u", &arangodb.QueryOptions{
		BindVars: map[string]interface{}{
			"name": userName,
		},
	}, nil)

	// SAFE: ValidateQuery with static string
	trx.ValidateQuery(ctx, "FOR u IN users FILTER u.age > 18 RETURN u")

	// SAFE: ExplainQuery with bind vars
	trx.ExplainQuery(ctx, "FOR u IN users FILTER u.name == @name RETURN u", map[string]interface{}{
		"name": userName,
	}, nil)
}

func transactionQueryInjectionWithControlFlow(db arangodb.Database, userName string, includeAge bool, userAge int) {
	ctx := context.Background()

	trx, _ := db.BeginTransaction(ctx, arangodb.TransactionCollections{
		Read: []string{"users"},
	}, &arangodb.BeginTransactionOptions{AllowImplicit: false})

	// UNSAFE: Concatenation in if block
	var query string
	if includeAge {
		query = "FOR u IN users FILTER u.name == '" + userName + "' AND u.age == " + fmt.Sprint(userAge) + " RETURN u"
	} else {
		query = "FOR u IN users FILTER u.name == '" + userName + "' RETURN u"
	}
	trx.Query(ctx, query, nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: Building query in loop with concatenation
	filters := []string{userName, "john", "jane"}
	for _, name := range filters {
		loopQuery := "FOR u IN users FILTER u.name == '" + name + "' RETURN u"
		trx.Query(ctx, loopQuery, nil) // want "query string uses concatenation instead of bind variables"
	}

	// SAFE: Using bind vars with control flow
	var safeQuery string
	var bindVars map[string]interface{}
	if includeAge {
		safeQuery = "FOR u IN users FILTER u.name == @name AND u.age == @age RETURN u"
		bindVars = map[string]interface{}{
			"name": userName,
			"age":  userAge,
		}
	} else {
		safeQuery = "FOR u IN users FILTER u.name == @name RETURN u"
		bindVars = map[string]interface{}{
			"name": userName,
		}
	}
	trx.Query(ctx, safeQuery, &arangodb.QueryOptions{BindVars: bindVars})

	// SAFE: Building query in loop with bind vars
	for _, name := range filters {
		trx.Query(ctx, "FOR u IN users FILTER u.name == @name RETURN u", &arangodb.QueryOptions{
			BindVars: map[string]interface{}{
				"name": name,
			},
		})
	}
}

func transactionQueryInjectionEdgeCases(db arangodb.Database) {
	ctx := context.Background()

	trx, _ := db.BeginTransaction(ctx, arangodb.TransactionCollections{
		Read: []string{"users"},
	}, &arangodb.BeginTransactionOptions{AllowImplicit: false})

	// SAFE: Empty options
	trx.Query(ctx, "FOR u IN users RETURN u", &arangodb.QueryOptions{})

	// SAFE: Just batch size option
	trx.Query(ctx, "FOR u IN users RETURN u", &arangodb.QueryOptions{
		BatchSize: 100,
	})

	// SAFE: Complex static query (no variable interpolation)
	staticComplexQuery := `
		FOR u IN users
		FILTER u.age > 18
		FILTER u.status == 'active'
		SORT u.name ASC
		LIMIT 10
		RETURN u
	`
	trx.Query(ctx, staticComplexQuery, nil)

	const staticConst = "FOR u IN users RETURN u"
	// SAFE: Using a const
	trx.Query(ctx, staticConst, nil)
}

func transactionQueryInjectionWithHelpers(db arangodb.Database, userName string) {
	ctx := context.Background()

	trx, _ := db.BeginTransaction(ctx, arangodb.TransactionCollections{
		Read: []string{"users"},
	}, &arangodb.BeginTransactionOptions{AllowImplicit: false})

	// SAFE: Using helper that returns bind vars (conservative - not flagged)
	query, bindVars := buildSafeQuery(userName)
	trx.Query(ctx, query, &arangodb.QueryOptions{BindVars: bindVars})

	// This call site won't be flagged because unsafeQuery is just an identifier
	// The analyzer can't trace back through the function call
	unsafeQuery := buildUnsafeQuery(userName)
	trx.Query(ctx, unsafeQuery, nil)

	// UNSAFE: But if we build the query with concatenation in the same function, it will be flagged
	localUnsafeQuery := "FOR u IN users FILTER u.name == '" + userName + "' RETURN u"
	trx.Query(ctx, localUnsafeQuery, nil) // want "query string uses concatenation instead of bind variables"
}

// Test with transaction stored in variable
func transactionInVariable(db arangodb.Database, userName string) {
	ctx := context.Background()

	trx, _ := db.BeginTransaction(ctx, arangodb.TransactionCollections{
		Read: []string{"users"},
	}, &arangodb.BeginTransactionOptions{AllowImplicit: false})

	// Store transaction in a variable
	myTrx := trx

	// UNSAFE: Query on stored transaction
	myTrx.Query(ctx, "FOR u IN users FILTER u.name == '"+userName+"' RETURN u", nil) // want "query string uses concatenation instead of bind variables"

	// SAFE: Using bind vars
	myTrx.Query(ctx, "FOR u IN users FILTER u.name == @name RETURN u", &arangodb.QueryOptions{
		BindVars: map[string]interface{}{
			"name": userName,
		},
	})
}
