package evaluator

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mcjcloud/taurine/ast"
	"github.com/mcjcloud/taurine/util"
)

func executeStatement(stmt ast.Statement, scope *Scope) error {
	switch t := stmt.(type) {
	case *ast.EtchStatement:
		return executeEtchStatement(t, scope)
	case *ast.ReadStatement:
		return executeReadStatement(t, scope)
	case *ast.ExpressionStatement:
		_, err := evaluateExpression(t.Expression, scope)
		return err
	case *ast.BlockStatement:
		return executeBlockStatement(t, scope)
	case *ast.IfStatement:
		return executeIfStatement(t, scope)
	case *ast.ForLoopStatement:
		return executeForStatement(t, scope)
	case *ast.WhileLoopStatement:
		return executeWhileStatement(t, scope)
	case *ast.ReturnStatement:
		return executeReturnStatement(t, scope)
	default:
		return fmt.Errorf("unkown statement %s", stmt)
	}
}

func executeBlockStatement(blockStmt *ast.BlockStatement, scope *Scope) error {
	subScope := NewScopeWithParent(scope)
	for _, s := range blockStmt.Statements {
		err := executeStatement(s, subScope)
		if err != nil {
			return err
		}
		// if a block exists within the current scope, the return value should propogate up
		if subScope.ReturnValue != nil {
			scope.ReturnValue = subScope.ReturnValue
			break
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
			idVal := scope.Get(idExp.Name)
			if idVal != nil {
				toEtch = append(toEtch, idVal.String())
			} else {
				toEtch = append(toEtch, "nil")
			}
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

func executeImportStatement(stmt *ast.ImportStatement, scope *Scope, tree *ast.Ast, g *util.ImportGraph) error {
	absPath := util.ResolveImport(filepath.Dir(tree.FilePath), stmt.Source)
	// check that the referenced ast has been evaluated
	var node *util.ImportNode
	if n, ok := g.Nodes[absPath]; !ok {
		return fmt.Errorf("could not find referenced file %s", absPath)
	} else if !n.Ast.Evaluated {
		evalErr := Evaluate(n.Ast, g)
		if evalErr != nil {
			return evalErr
		}
		n.Ast.Evaluated = true
		node = n
	} else {
		node = n
	}

	// imported values should now exist in the Ast exports
	// add all the evaluated exports to the scope
	for _, id := range stmt.Imports {
		if exp, ok := node.Ast.Exports[id.Name]; !ok {
			return fmt.Errorf("symbol '%s' is not exported from %s", id.Name, absPath)
		} else {
			scope.Set(id.Name, exp)
		}
	}
	return nil
}

func executeExportStatement(stmt *ast.ExportStatement, scope *Scope, tree *ast.Ast) error {
	// evaluate the expression value
	val, err := evaluateExpression(stmt.Value, scope)
	if err != nil {
		return err
	}
	tree.Exports[stmt.Identifier.Name] = val
	return nil
}

func executeIfStatement(ifStmt *ast.IfStatement, scope *Scope) error {
	exp, err := evaluateExpression(ifStmt.Condition, scope)
	if err != nil {
		return err
	}
	if boolExp, ok := exp.(*ast.BooleanLiteral); ok {
		if boolExp.Value {
			if err := executeStatement(ifStmt.Statement, scope); err != nil {
				return err
			}
		} else if ifStmt.ElseIf != nil {
			if err := executeStatement(ifStmt.ElseIf, scope); err != nil {
				return err
			}
		}
		return nil
	}
	return errors.New("if expression must evaluate to boolean")
}

func executeForStatement(forStmt *ast.ForLoopStatement, scope *Scope) error {
	// evaluate the iterator
	arrExp, err := evaluateExpression(forStmt.Iterator, scope)
	if err != nil {
		return err
	}
	var arr *ast.ArrayExpression
	if a, ok := arrExp.(*ast.ArrayExpression); ok {
		arr = a
	} else if s, ok := arrExp.(*ast.StringLiteral); ok {
		chars := make([]ast.Expression, len(s.Value))
		for i, c := range s.Value {
			chars[i] = &ast.StringLiteral{Value: string(c)}
		}
		arr = &ast.ArrayExpression{Expressions: chars}
	} else {
		return fmt.Errorf("expected array or string iterator but found %s", arrExp)
	}

	// loop through the array
	for i := 0; i < len(arr.Expressions); i += forStmt.Step {
		control, err := evaluateExpression(arr.Expressions[i], scope)
		if err != nil {
			return err
		}
		forScope := NewScopeWithParent(scope)
		forScope.Set(forStmt.Control.Name, control)
		if err := executeStatement(forStmt.Statement, forScope); err != nil {
			return err
		}
		if forScope.ReturnValue != nil {
			scope.ReturnValue = forScope.ReturnValue
			break
		}
	}
	return nil
}

func executeWhileStatement(whileStmt *ast.WhileLoopStatement, scope *Scope) error {
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
			// if there is a return value, the loop should end
			if subScope.ReturnValue != nil {
				scope.ReturnValue = subScope.ReturnValue
				break
			}
		}
	} else {
		return errors.New("while expression must evaluate to boolean")
	}
	return nil
}

func executeReturnStatement(rtnStmt *ast.ReturnStatement, scope *Scope) error {
	exp, err := evaluateExpression(rtnStmt.Value, scope)
	if err != nil {
		return err
	}
	scope.ReturnValue = exp
	return nil
}
