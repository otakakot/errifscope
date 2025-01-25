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
		switch n := node.(type) {
		case *ast.IfStmt:
			binaryExpr, ok := n.Cond.(*ast.BinaryExpr)
			if !ok {
				return
			}

			assignStmt, ok := previousStmt.(*ast.AssignStmt)
			if !ok {
				return
			}

			// for error variable check
			errFlag := false

			// for value variable check
			valFlag := false

			var errExpr ast.Expr

			for _, lh := range assignStmt.Lhs {
				typeof := pass.TypesInfo.TypeOf(lh)
				if typeof == nil {
					continue
				}

				if typeof.String() != "error" {
					lhIdent, ok := lh.(*ast.Ident)
					if !ok {
						continue
					}

					if lhIdent.Name != "_" {
						valFlag = true

						continue
					}
				}

				errFlag = true

				errExpr = lh
			}

			if errFlag && valFlag {
				return
			}

			typeof := pass.TypesInfo.TypeOf(errExpr)
			if typeof == nil {
				return
			}

			if typeof.String() != "error" {
				return
			}

			errIdent, ok := errExpr.(*ast.Ident)
			if !ok {
				return
			}

			ifErrIdent, ok := binaryExpr.X.(*ast.Ident)
			if !ok {
				return
			}

			if errIdent.Name != ifErrIdent.Name {
				return
			}

			if binaryExpr.Op.String() != "!=" {
				return
			}

			yIdent, ok := binaryExpr.Y.(*ast.Ident)
			if !ok {
				return
			}

			if yIdent.Obj != nil {
				return
			}

			pass.Reportf(node.Pos(), msg, errIdent.Name)
		}
		previousStmt = node
	})

	return nil, nil
}
