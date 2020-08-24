package evaluator

import (
	"fmt"
	"strings"

	"github.com/mcjcloud/taurine/ast"
)

// Evaluate evaluates the code and does stuff
func Evaluate(block *ast.BlockStatement) error {
	scope := &Scope{}
	for _, stmt := range block.Statements {
		// do the statement
		err := executeStatement(stmt, scope)
		if err != nil {
			return err
		}
	}
	return nil
}

func executeStatement(stmt ast.Statement, scope *Scope) error {
	if etchStmt, ok := stmt.(*ast.EtchStatement); ok {
		if err := executeEtchStatement(etchStmt, scope); err != nil {
			return err
		}
	} else if declStmt, ok := stmt.(*ast.VariableDecleration); ok {
		if err := executeVariableDecleration(declStmt, scope); err != nil {
			return err
		}
	}
	return nil
}

func executeEtchStatement(stmt *ast.EtchStatement, scope *Scope) error {
	var toEtch []string
	for _, exp := range stmt.Expressions {
		if numExp, ok := exp.(*ast.NumberLiteral); ok {
			toEtch = append(toEtch, numExp.String())
		} else if strExp, ok := exp.(*ast.StringLiteral); ok {
			toEtch = append(toEtch, strExp.String())
		} else if idExp, ok := exp.(*ast.Identifier); ok {
			toEtch = append(toEtch, scope.Get(idExp.Name).String())
		}
	}
	fmt.Println(strings.Join(toEtch, " "))
	return nil
}

func executeVariableDecleration(stmt *ast.VariableDecleration, scope *Scope) error {
	val, err := evaluateExpression(stmt.Value, scope)
	if err != nil {
		return err
	}
	scope.Set(stmt.Symbol, val)
	return nil
}

func evaluateExpression(exp ast.Expression, scope *Scope) (ast.Expression, error) {
	if op, ok := exp.(*ast.OperationExpression); ok {
		return evaluateOperation(op, scope)
	} else if id, ok := exp.(*ast.Identifier); ok {
		return scope.Get(id.Name), nil
	}
	return exp, nil
}
