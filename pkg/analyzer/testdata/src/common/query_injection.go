package common

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver/v2/arangodb"
)

func queryInjectionExamples(db arangodb.Database, userName string, userAge int) {
	ctx := context.Background()

	// UNSAFE: Direct string concatenation with +
	db.Query(ctx, "FOR u IN users FILTER u.name == '"+userName+"' RETURN u", nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: String concatenation across multiple lines
	query := "FOR u IN users FILTER u.name == '" + userName + "' RETURN u"
	db.Query(ctx, query, nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: Multiple concatenations
	db.Query(ctx, "FOR u IN users FILTER u.name == '"+userName+"' AND u.age == "+fmt.Sprint(userAge)+" RETURN u", nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: Concatenation with compound assignment +=
	q := "FOR u IN users"
	q += " FILTER u.name == '" + userName + "'"
	q += " RETURN u"
	db.Query(ctx, q, nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: fmt.Sprintf
	sprintfQuery := fmt.Sprintf("FOR u IN users FILTER u.name == '%s' RETURN u", userName)
	db.Query(ctx, sprintfQuery, nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: QueryBatch with concatenation
	db.QueryBatch(ctx, "FOR u IN users FILTER u.name == '"+userName+"' RETURN u", nil, nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: ValidateQuery with concatenation
	db.ValidateQuery(ctx, "FOR u IN users FILTER u.name == '"+userName+"' RETURN u") // want "query string uses concatenation instead of bind variables"

	// UNSAFE: ExplainQuery with concatenation (note: ExplainQuery has bindVars param but query could still be built with concat)
	db.ExplainQuery(ctx, "FOR u IN users FILTER u.name == '"+userName+"' RETURN u", nil, nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: Concatenation in variable initialization
	var concatenatedQuery = "FOR u IN users FILTER u.name == '" + userName + "' RETURN u"
	db.Query(ctx, concatenatedQuery, nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: Building query step by step with concatenation
	baseQuery := "FOR u IN users"
	filterPart := " FILTER u.name == '" + userName + "'"
	finalQuery := baseQuery + filterPart + " RETURN u"
	db.Query(ctx, finalQuery, nil) // want "query string uses concatenation instead of bind variables"

	// SAFE: Using bind variables properly
	db.Query(ctx, "FOR u IN users FILTER u.name == @name RETURN u", &arangodb.QueryOptions{
		BindVars: map[string]interface{}{
			"name": userName,
		},
	})

	// SAFE: Using bind variables with multiple parameters
	db.Query(ctx, "FOR u IN users FILTER u.name == @name AND u.age == @age RETURN u", &arangodb.QueryOptions{
		BindVars: map[string]interface{}{
			"name": userName,
			"age":  userAge,
		},
	})

	// SAFE: Static query string (no variables)
	db.Query(ctx, "FOR u IN users RETURN u", nil)

	// SAFE: Static query with options but no variables
	db.Query(ctx, "FOR u IN users FILTER u.age > 18 RETURN u", &arangodb.QueryOptions{
		Count: true,
	})

	// SAFE: Concatenation of only static strings (no variables) - this is OK
	staticQuery := "FOR u IN users" + " FILTER u.age > 18" + " RETURN u"
	db.Query(ctx, staticQuery, nil)

	// SAFE: QueryBatch with bind vars
	db.QueryBatch(ctx, "FOR u IN users FILTER u.name == @name RETURN u", &arangodb.QueryOptions{
		BindVars: map[string]interface{}{
			"name": userName,
		},
	}, nil)

	// SAFE: ValidateQuery with static string
	db.ValidateQuery(ctx, "FOR u IN users FILTER u.age > 18 RETURN u")

	// SAFE: ExplainQuery with bind vars (proper usage)
	db.ExplainQuery(ctx, "FOR u IN users FILTER u.name == @name RETURN u", map[string]interface{}{
		"name": userName,
	}, nil)

	// UNSAFE: Concatenation in a more complex expression
	prefix := "FOR u IN users FILTER "
	condition := "u.name == '" + userName + "'"
	suffix := " RETURN u"
	complexQuery := prefix + condition + suffix
	db.Query(ctx, complexQuery, nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: Using fmt.Sprintf with multiple args
	multiSprintfQuery := fmt.Sprintf("FOR u IN users FILTER u.name == '%s' AND u.age == %d RETURN u", userName, userAge)
	db.Query(ctx, multiSprintfQuery, nil) // want "query string uses concatenation instead of bind variables"
}

func queryInjectionWithControlFlow(db arangodb.Database, userName string, includeAge bool, userAge int) {
	ctx := context.Background()

	// UNSAFE: Concatenation in if block
	var query string
	if includeAge {
		query = "FOR u IN users FILTER u.name == '" + userName + "' AND u.age == " + fmt.Sprint(userAge) + " RETURN u"
	} else {
		query = "FOR u IN users FILTER u.name == '" + userName + "' RETURN u"
	}
	db.Query(ctx, query, nil) // want "query string uses concatenation instead of bind variables"

	// UNSAFE: Building query in loop with concatenation
	filters := []string{userName, "john", "jane"}
	for _, name := range filters {
		loopQuery := "FOR u IN users FILTER u.name == '" + name + "' RETURN u"
		db.Query(ctx, loopQuery, nil) // want "query string uses concatenation instead of bind variables"
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
	db.Query(ctx, safeQuery, &arangodb.QueryOptions{BindVars: bindVars})

	// SAFE: Building query in loop with bind vars
	for _, name := range filters {
		db.Query(ctx, "FOR u IN users FILTER u.name == @name RETURN u", &arangodb.QueryOptions{
			BindVars: map[string]interface{}{
				"name": name,
			},
		})
	}
}

func queryInjectionEdgeCases(db arangodb.Database) {
	ctx := context.Background()

	// SAFE: Empty options
	db.Query(ctx, "FOR u IN users RETURN u", &arangodb.QueryOptions{})

	// SAFE: Just batch size option
	db.Query(ctx, "FOR u IN users RETURN u", &arangodb.QueryOptions{
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
	db.Query(ctx, staticComplexQuery, nil)

	const staticConst = "FOR u IN users RETURN u"
	// SAFE: Using a const
	db.Query(ctx, staticConst, nil)
}

// Helper functions that return queries
// Note: The analyzer is conservative and doesn't trace through function calls.
// It only flags concatenation at the Query call site or in prior statements.
func buildSafeQuery(name string) (string, map[string]interface{}) {
	return "FOR u IN users FILTER u.name == @name RETURN u", map[string]interface{}{
		"name": name,
	}
}

func buildUnsafeQuery(name string) string {
	// This concatenation won't be flagged because it's in a return statement
	// The analyzer only checks db.Query() call sites
	return "FOR u IN users FILTER u.name == '" + name + "' RETURN u"
}

func queryInjectionWithHelpers(db arangodb.Database, userName string) {
	ctx := context.Background()

	// SAFE: Using helper that returns bind vars
	query, bindVars := buildSafeQuery(userName)
	db.Query(ctx, query, &arangodb.QueryOptions{BindVars: bindVars})

	// This call site won't be flagged because unsafeQuery is just an identifier
	// The analyzer can't trace back through the function call
	unsafeQuery := buildUnsafeQuery(userName)
	db.Query(ctx, unsafeQuery, nil)

	// UNSAFE: But if we build the query with concatenation in the same function, it will be flagged
	localUnsafeQuery := "FOR u IN users FILTER u.name == '" + userName + "' RETURN u"
	db.Query(ctx, localUnsafeQuery, nil) // want "query string uses concatenation instead of bind variables"
}
