package errifscope

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "errifscope",
	Doc:  "errifscope is linter to find if block that can encapsulate the scope of error variable.",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

const msg = "%s can be scoped with if block"

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.AssignStmt)(nil),
		(*ast.IfStmt)(nil),
	}

	var previousStmt ast.Node

	inspect.Preorder(nodeFilter, func(node ast.Node) {
		if ifStmt, ok := node.(*ast.IfStmt); ok {
			processIfStatement(pass, ifStmt, previousStmt)
		}
		previousStmt = node
	})

	return nil, nil
}

func processIfStatement(
	pass *analysis.Pass,
	ifStmt *ast.IfStmt,
	previousStmt ast.Node,
) {
	binaryExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		return
	}

	assignStmt, ok := previousStmt.(*ast.AssignStmt)
	if !ok {
		return
	}

	errExpr, isScopable := analyzeAssignment(pass, assignStmt)
	if !isScopable {
		return
	}

	if isErrorCheckValid(pass, binaryExpr, errExpr) {
		errIdent := errExpr.(*ast.Ident)
		pass.Reportf(ifStmt.Pos(), msg, errIdent.Name)
	}
}

func analyzeAssignment(
	pass *analysis.Pass,
	assignStmt *ast.AssignStmt,
) (ast.Expr, bool) {
	var errExpr ast.Expr

	hasErrorVar, hasValueVar := false, false

	for _, lhs := range assignStmt.Lhs {
		typeof := pass.TypesInfo.TypeOf(lhs)
		if typeof == nil {
			continue
		}

		if typeof.String() == "error" {
			hasErrorVar = true
			errExpr = lhs
		} else if ident, ok := lhs.(*ast.Ident); ok && ident.Name != "_" {
			hasValueVar = true
		}
	}

	return errExpr, hasErrorVar && !hasValueVar
}

func isErrorCheckValid(
	pass *analysis.Pass,
	binaryExpr *ast.BinaryExpr,
	errExpr ast.Expr,
) bool {
	if pass.TypesInfo.TypeOf(errExpr).String() != "error" {
		return false
	}

	errIdent, ok := errExpr.(*ast.Ident)
	if !ok {
		return false
	}

	ifIdent, ok := binaryExpr.X.(*ast.Ident)
	if !ok || errIdent.Name != ifIdent.Name {
		return false
	}

	if binaryExpr.Op.String() != "!=" {
		return false
	}

	yIdent, ok := binaryExpr.Y.(*ast.Ident)

	return ok && yIdent.Obj == nil
}
