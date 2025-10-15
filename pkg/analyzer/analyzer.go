// Package analyzer contains tools for analyzing arangodb usage.
// It focuses on github.com/arangodb/go-driver/v2.
package analyzer

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NewAnalyzer returns an arangolint analyzer.
func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "arangolint",
		Doc:      "opinionated best practices for arangodb client",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

var (
	errUnknown         = errors.New("unknown node type")
	errInvalidAnalysis = errors.New("invalid analysis")
)

const missingAllowImplicitOptionMsg = "missing AllowImplicit option"

func run(pass *analysis.Pass) (interface{}, error) {
	inspctr, typeValid := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !typeValid {
		return nil, errInvalidAnalysis
	}

	var stack []ast.Node

	inspctr.Nodes(nil, func(node ast.Node, push bool) (proceed bool) {
		// pop
		if !push {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]

				return true
			}
		}

		// push
		stack = append(stack, node)

		if call, isCall := node.(*ast.CallExpr); isCall {
			handleBeginTransactionCall(call, pass, stack)
		}

		return true
	})

	return nil, nil //nolint:nilnil
}

// handleBeginTransactionCall encapsulates the logic for validating BeginTransaction calls
// to keep the cognitive complexity of run() low.
func handleBeginTransactionCall(call *ast.CallExpr, pass *analysis.Pass, stack []ast.Node) {
	if !isBeginTransaction(call, pass) {
		return
	}

	diag := analysis.Diagnostic{
		Pos:     call.Pos(),
		Message: missingAllowImplicitOptionMsg,
	}

	switch typedArg := call.Args[2].(type) {
	case *ast.Ident:
		if typedArg.Name == "nil" {
			pass.Report(diag)

			return
		}

		if hasAllowImplicitForIdent(typedArg, pass, stack, call.Pos()) {
			return
		}

		pass.Report(diag)
	case *ast.UnaryExpr:
		// &literal or &ident
		elts, err := getElts(typedArg.X)
		if err == nil {
			if !eltsHasAllowImplicit(elts) {
				pass.Report(diag)
			}

			return
		}

		// not a literal, try &ident
		if id, ok := typedArg.X.(*ast.Ident); ok {
			if hasAllowImplicitForIdent(id, pass, stack, call.Pos()) {
				return
			}

			pass.Report(diag)
		}
	}
}

func isBeginTransaction(call *ast.CallExpr, pass *analysis.Pass) bool {
	selExpr, isSelector := call.Fun.(*ast.SelectorExpr)
	if !isSelector {
		return false
	}

	xType := pass.TypesInfo.TypeOf(selExpr.X)
	if xType == nil {
		return false
	}

	const arangoStruct = "github.com/arangodb/go-driver/v2/arangodb.Database"

	if !strings.HasSuffix(xType.String(), arangoStruct) ||
		selExpr.Sel.Name != "BeginTransaction" {
		return false
	}

	const expectedArgsCount = 3

	return len(call.Args) == expectedArgsCount
}

func getElts(node ast.Node) ([]ast.Expr, error) {
	switch typedNode := node.(type) {
	case *ast.CompositeLit:
		return typedNode.Elts, nil
	default:
		return nil, errUnknown
	}
}

func eltsHasAllowImplicit(elts []ast.Expr) bool {
	for _, elt := range elts {
		if eltIsAllowImplicit(elt) {
			return true
		}
	}

	return false
}

func eltIsAllowImplicit(expr ast.Expr) bool {
	switch typedNode := expr.(type) {
	case *ast.KeyValueExpr:
		ident, ok := typedNode.Key.(*ast.Ident)
		if !ok {
			return false
		}

		return ident.Name == "AllowImplicit"
	default:
		return false
	}
}

// hasAllowImplicitForIdent checks whether the given identifier (variable or pointer to options)
// has the AllowImplicit option explicitly set before the call position within the nearest enclosing block.
func hasAllowImplicitForIdent(
	id *ast.Ident,
	pass *analysis.Pass,
	stack []ast.Node,
	callPos token.Pos,
) bool {
	obj := pass.TypesInfo.ObjectOf(id)
	if obj == nil {
		return false
	}

	blk := nearestEnclosingBlock(stack)
	if blk == nil {
		return false
	}

	// scan statements in order until the call position
	for _, stmt := range blk.List {
		if stmt == nil {
			continue
		}

		if stmt.Pos() >= callPos {
			break
		}

		// explicit assignment like: options.AllowImplicit = ...
		if hasAllowImplicitAssignForObj(stmt, obj, pass) {
			return true
		}

		// initialization via short/regular assignment: options := {AllowImplicit: ...} or options = {AllowImplicit: ...}
		if as, ok := stmt.(*ast.AssignStmt); ok {
			if initHasAllowImplicitForObj(as, obj, pass) {
				return true
			}
		}

		// initialization via var declaration: var options = {AllowImplicit: ...} or var optns = &{AllowImplicit: ...}
		if declInitHasAllowImplicitForObj(stmt, obj, pass) {
			return true
		}
	}

	return false
}

func nearestEnclosingBlock(stack []ast.Node) *ast.BlockStmt {
	for i := len(stack) - 1; i >= 0; i-- {
		if blk, ok := stack[i].(*ast.BlockStmt); ok {
			return blk
		}
	}

	return nil
}

func hasAllowImplicitAssignForObj(stmt ast.Stmt, obj types.Object, pass *analysis.Pass) bool {
	as, isAssignStmt := stmt.(*ast.AssignStmt)
	if !isAssignStmt {
		return false
	}

	for _, lhs := range as.Lhs {
		sel, isSelectorExpr := lhs.(*ast.SelectorExpr)
		if !isSelectorExpr {
			continue
		}

		if sel.Sel == nil || sel.Sel.Name != "AllowImplicit" {
			continue
		}

		ident, isIdent := sel.X.(*ast.Ident)
		if !isIdent {
			continue
		}

		if pass.TypesInfo.ObjectOf(ident) == obj {
			return true
		}
	}

	return false
}

func initHasAllowImplicitForObj(
	assign *ast.AssignStmt,
	obj types.Object,
	pass *analysis.Pass,
) bool {
	// find the RHS corresponding to our obj
	for lhsIndex, lhs := range assign.Lhs {
		id, ok := lhs.(*ast.Ident)
		if !ok {
			continue
		}

		if pass.TypesInfo.ObjectOf(id) != obj {
			continue
		}

		var rhs ast.Expr

		switch {
		case len(assign.Rhs) == len(assign.Lhs):
			rhs = assign.Rhs[lhsIndex]
		case len(assign.Rhs) == 1:
			rhs = assign.Rhs[0]
		default:
			continue
		}

		// allow either &CompositeLit or CompositeLit
		if ue, ok := rhs.(*ast.UnaryExpr); ok {
			elts, err := getElts(ue.X)
			if err == nil {
				return eltsHasAllowImplicit(elts)
			}
		}

		if cl, ok := rhs.(*ast.CompositeLit); ok {
			return eltsHasAllowImplicit(cl.Elts)
		}
	}

	return false
}

func declInitHasAllowImplicitForObj(stmt ast.Stmt, obj types.Object, pass *analysis.Pass) bool {
	declStmt, isDeclStmt := stmt.(*ast.DeclStmt)
	if !isDeclStmt {
		return false
	}

	genDecl, isGenDecl := declStmt.Decl.(*ast.GenDecl)
	if !isGenDecl || genDecl.Tok != token.VAR {
		return false
	}

	for _, spec := range genDecl.Specs {
		valueSpec, isValueSpec := spec.(*ast.ValueSpec)
		if !isValueSpec {
			continue
		}

		if valueSpecHasAllowImplicitForObj(valueSpec, obj, pass) {
			return true
		}
	}

	return false
}

func valueSpecHasAllowImplicitForObj(
	valueSpec *ast.ValueSpec,
	obj types.Object,
	pass *analysis.Pass,
) bool {
	// find the index corresponding to our obj
	targetIndex := -1

	for i, name := range valueSpec.Names {
		if pass.TypesInfo.ObjectOf(name) == obj {
			targetIndex = i

			break
		}
	}

	if targetIndex == -1 {
		return false
	}

	// pick the value expression for this name
	var value ast.Expr

	switch {
	case targetIndex < len(valueSpec.Values):
		value = valueSpec.Values[targetIndex]
	case len(valueSpec.Values) == 1:
		value = valueSpec.Values[0]
	default:
		return false
	}

	// allow either &CompositeLit or CompositeLit
	if ue, isUnary := value.(*ast.UnaryExpr); isUnary {
		elts, err := getElts(ue.X)
		if err == nil {
			return eltsHasAllowImplicit(elts)
		}
	}

	if cl, isComposite := value.(*ast.CompositeLit); isComposite {
		return eltsHasAllowImplicit(cl.Elts)
	}

	return false
}
