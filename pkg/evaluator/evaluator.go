package evaluator

import (
	"errors"
	"fmt"

	"github.com/mcjcloud/taurine/pkg/ast"
	"github.com/mcjcloud/taurine/pkg/util"
)

// Evaluate evaluates the code and does stuff
func Evaluate(tree *ast.Ast, importGraph *util.ImportGraph) error {
	// check that the ast has a blockstatement
	var block *ast.BlockStatement
	if b, ok := tree.Statement.(*ast.BlockStatement); !ok {
		return errors.New("ast must contain block statement")
	} else {
		block = b
	}

	// execute block statements
	scope := NewScope()
	for _, stmt := range block.Statements {
		if importStmt, ok := stmt.(*ast.ImportStatement); ok {
			if err := executeImportStatement(importStmt, scope, tree, importGraph); err != nil {
				return err
			}
		} else if exportStmt, ok := stmt.(*ast.ExportStatement); ok {
			if err := executeExportStatement(exportStmt, scope, tree); err != nil {
				return err
			}
		} else {
			// do the statement
			err := executeStatement(stmt, scope)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func conformDataType(dType ast.Symbol, exp ast.Expression) (ast.Expression, error) {
	switch dType {
	case ast.NUM:
		if _, nok := exp.(*ast.NumberLiteral); nok {
			return exp, nil
		}
		if iExp, iok := exp.(*ast.IntegerLiteral); iok {
			return &ast.NumberLiteral{Value: float64(iExp.Value.Int64())}, nil
		}
	case ast.INT:
		if _, iok := exp.(*ast.IntegerLiteral); iok {
			return exp, nil
		}
	case ast.STR:
		if _, sok := exp.(*ast.StringLiteral); sok {
			return exp, nil
		}
	case ast.BOOL:
		if _, iok := exp.(*ast.BooleanLiteral); iok {
			return exp, nil
		}
	case ast.ARR:
		if _, iok := exp.(*ast.ArrayExpression); iok {
			return exp, nil
		}
	case ast.FUNC:
		if _, fok := exp.(*ast.FunctionLiteral); fok {
			return exp, nil
		}
		if _, fok := exp.(*ScopedFunction); fok {
			return exp, nil
		}
	}
	return nil, fmt.Errorf("%s is not of type %s", exp, dType)
}
