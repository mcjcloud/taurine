package evaluator

import (
	"errors"
	"fmt"

	"github.com/mcjcloud/taurine/ast"
)

func evaluateOperation(op *ast.OperationExpression, scope *Scope) (ast.Expression, error) {
	left, err := evaluateExpression(op.LeftExpression, scope)
	if err != nil {
		return nil, err
	}
	right, err := evaluateExpression(op.RightExpression, scope)
	if err != nil {
		return nil, err
	}
	if op.Operator == "+" {
		if leftNum, ok := left.(*ast.NumberLiteral); ok {
			if rightNum, ok := right.(*ast.NumberLiteral); ok {
				return &ast.NumberLiteral{Value: leftNum.Value + rightNum.Value}, nil
			} else if rightStr, ok := right.(*ast.StringLiteral); ok {
				return &ast.StringLiteral{Value: fmt.Sprintf("%f%s", leftNum.Value, rightStr.Value)}, nil
			}
		} else if leftStr, ok := left.(*ast.StringLiteral); ok {
			if rightNum, ok := right.(*ast.NumberLiteral); ok {
				return &ast.StringLiteral{Value: fmt.Sprintf("%s%f", leftStr.Value, rightNum.Value)}, nil
			} else if rightStr, ok := right.(*ast.StringLiteral); ok {
				return &ast.StringLiteral{Value: fmt.Sprintf("%s%s", leftStr.Value, rightStr.Value)}, nil
			}
		}
		return nil, errors.New("'+' operator is not applicable to arguments")
	} else if op.Operator == "-" {
		if leftNum, ok := left.(*ast.NumberLiteral); ok {
			if rightNum, ok := right.(*ast.NumberLiteral); ok {
				return &ast.NumberLiteral{Value: leftNum.Value - rightNum.Value}, nil
			}
		}
		return nil, errors.New("'-' operator only applies to type num")
	} else if op.Operator == "/" {
		if leftNum, ok := left.(*ast.NumberLiteral); ok {
			if rightNum, ok := right.(*ast.NumberLiteral); ok {
				return &ast.NumberLiteral{Value: leftNum.Value / rightNum.Value}, nil
			}
		}
		return nil, errors.New("'/' operator only applies to type num")
	} else if op.Operator == "*" {
		if leftNum, ok := left.(*ast.NumberLiteral); ok {
			if rightNum, ok := right.(*ast.NumberLiteral); ok {
				return &ast.NumberLiteral{Value: leftNum.Value * rightNum.Value}, nil
			}
		}
		return nil, errors.New("'*' operator only applies to type num")
	}
	return nil, errors.New("unrecognized operator")
}
