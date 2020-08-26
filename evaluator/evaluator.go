package evaluator

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mcjcloud/taurine/ast"
)

// Evaluate evaluates the code and does stuff
func Evaluate(block *ast.BlockStatement) error {
	scope := NewScope()
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
	} else if readStmt, ok := stmt.(*ast.ReadStatement); ok {
		if err := executeReadStatement(readStmt, scope); err != nil {
			return err
		}
	} else if declStmt, ok := stmt.(*ast.VariableDecleration); ok {
		if err := executeVariableDecleration(declStmt, scope); err != nil {
			return err
		}
	} else if expStmt, ok := stmt.(*ast.ExpressionStatement); ok {
		_, err := evaluateExpression(expStmt.Expression, scope)
		return err
	} else if blockStmt, ok := stmt.(*ast.BlockStatement); ok {
		subScope := NewScopeWithParent(scope)
		for _, s := range blockStmt.Statements {
			err := executeStatement(s, subScope)
			if err != nil {
				return err
			}
		}
	} else if whileStmt, ok := stmt.(*ast.WhileLoopStatement); ok {
		exp, err := evaluateExpression(whileStmt.Condition, scope)
		if err != nil {
			return err
		}
		if boolExp, ok := exp.(*ast.BooleanLiteral); ok {
			subScope := NewScopeWithParent(scope)
			for boolExp.Value {
				err := executeStatement(whileStmt.Statement, subScope)
				if err != nil {
					return err
				}
				exp, err = evaluateExpression(whileStmt.Condition, subScope)
				if err != nil {
					return err
				}
				boolExp, ok = exp.(*ast.BooleanLiteral)
				if !ok {
					return errors.New("while expression is no longer boolean")
				}
			}
		} else {
			return errors.New("while expression must evaluate to boolean")
		}
	} else {
		return errors.New("unrecognized statement")
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
		} else {
			expEval, err := evaluateExpression(exp, scope)
			if err != nil {
				return err
			}
			toEtch = append(toEtch, expEval.String())
		}
	}
	fmt.Println(strings.Join(toEtch, " "))
	return nil
}

func executeReadStatement(stmt *ast.ReadStatement, scope *Scope) error {
	if stmt.Prompt != nil {
		fmt.Printf("%s", stmt.Prompt)
	}
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		scope.Set(stmt.Identifier.Name, &ast.StringLiteral{Value: scanner.Text()})
	} else {
		return errors.New("error reading input")
	}
	return nil
}

func executeVariableDecleration(stmt *ast.VariableDecleration, scope *Scope) error {
	val, err := evaluateExpression(stmt.Value, scope)
	if err != nil {
		return err
	}
	if scope.Variables[stmt.Symbol] != nil || (scope.Parent != nil && scope.Parent.Get(stmt.Symbol) != nil) {
		return fmt.Errorf("variable '%s' already exists", stmt.Symbol)
	}
	scope.Set(stmt.Symbol, val)
	return nil
}

func evaluateExpression(exp ast.Expression, scope *Scope) (ast.Expression, error) {
	if op, ok := exp.(*ast.OperationExpression); ok {
		return evaluateOperation(op, scope)
	} else if id, ok := exp.(*ast.Identifier); ok {
		return scope.Get(id.Name), nil
	} else if asn, ok := exp.(*ast.AssignmentExpression); ok {
		// make sure the identifier exists
		if scope.Get(asn.Identifier.Name) == nil {
			return nil, fmt.Errorf("'%s' was not declared", asn.Identifier.Name)
		}
		val, err := evaluateExpression(asn.Value, scope)
		if err != nil {
			return nil, err
		}
		// update the scope and return the evaluated value
		scope.Set(asn.Identifier.Name, val)
		return val, nil
	}
	return exp, nil
}
