package evaluator

import (
	"fmt"

	"github.com/mcjcloud/taurine/pkg/ast"
)

func equalEqual(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	left, right, err := evaluateOperands(leftExp, rightExp, scope)
	if err != nil {
		return nil, err
	}

	leftNum, lok := left.(*ast.NumberLiteral)
	if lok {
		rightVal, err := conformDataType(ast.NUM, right)
		if err != nil {
			return nil, err
		}
		rightNum, rok := rightVal.(*ast.NumberLiteral)
		if rok {
			return &ast.BooleanLiteral{Value: leftNum.Value == rightNum.Value}, nil
		}
	}

	leftInt, lok := left.(*ast.IntegerLiteral)
	rightInt, rok := right.(*ast.IntegerLiteral)
	if lok && rok {
		return &ast.BooleanLiteral{Value: leftInt.Value.Cmp(rightInt.Value) == 0}, nil
	}

	leftStr, lok := left.(*ast.StringLiteral)
	rightStr, rok := right.(*ast.StringLiteral)
	if lok && rok {
		return &ast.BooleanLiteral{Value: leftStr.Value == rightStr.Value}, nil
	}

	leftBool, lok := left.(*ast.BooleanLiteral)
	rightBool, rok := right.(*ast.BooleanLiteral)
	if lok && rok {
		return &ast.BooleanLiteral{Value: leftBool.Value == rightBool.Value}, nil
	}

	return nil, fmt.Errorf("'==' cannot be applied to '%s' and '%s'", leftExp, rightExp)
}

func notEqual(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	left, right, err := evaluateOperands(leftExp, rightExp, scope)
	if err != nil {
		return nil, err
	}

	leftNum, lok := left.(*ast.NumberLiteral)
	if lok {
		rightVal, err := conformDataType(ast.NUM, right)
		if err != nil {
			return nil, err
		}
		rightNum, rok := rightVal.(*ast.NumberLiteral)
		if rok {
			return &ast.BooleanLiteral{Value: leftNum.Value != rightNum.Value}, nil
		}
	}

	leftInt, lok := left.(*ast.IntegerLiteral)
	rightInt, rok := right.(*ast.IntegerLiteral)
	if lok && rok {
		return &ast.BooleanLiteral{Value: leftInt.Value.Cmp(rightInt.Value) != 0}, nil
	}

	leftStr, lok := left.(*ast.StringLiteral)
	rightStr, rok := right.(*ast.StringLiteral)
	if lok && rok {
		return &ast.BooleanLiteral{Value: leftStr.Value != rightStr.Value}, nil
	}

	leftBool, lok := left.(*ast.BooleanLiteral)
	rightBool, rok := right.(*ast.BooleanLiteral)
	if lok && rok {
		return &ast.BooleanLiteral{Value: leftBool.Value != rightBool.Value}, nil
	}

	return nil, fmt.Errorf("'!=' cannot be applied to '%s' and '%s'", leftExp, rightExp)
}

func lessThan(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	left, right, err := evaluateOperands(leftExp, rightExp, scope)
	if err != nil {
		return nil, err
	}

	if leftNum, ok := left.(*ast.NumberLiteral); ok {
		rightVal, err := conformDataType(ast.NUM, right)
		if err != nil {
			return nil, err
		}
		if rightNum, ok := rightVal.(*ast.NumberLiteral); ok {
			return &ast.BooleanLiteral{Value: leftNum.Value < rightNum.Value}, nil
		}
	} else if leftInt, ok := left.(*ast.IntegerLiteral); ok {
		if rightInt, ok := right.(*ast.IntegerLiteral); ok {
			return &ast.BooleanLiteral{Value: leftInt.Value.Cmp(rightInt.Value) < 0}, nil
		}
	}
	return nil, fmt.Errorf("'<' cannot be applied to '%s' and '%s'", leftExp, rightExp)
}

func lessEqual(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	left, right, err := evaluateOperands(leftExp, rightExp, scope)
	if err != nil {
		return nil, err
	}

	if leftNum, ok := left.(*ast.NumberLiteral); ok {
		rightVal, err := conformDataType(ast.NUM, right)
		if err != nil {
			return nil, err
		}
		if rightNum, ok := rightVal.(*ast.NumberLiteral); ok {
			return &ast.BooleanLiteral{Value: leftNum.Value <= rightNum.Value}, nil
		}
	} else if leftInt, ok := left.(*ast.IntegerLiteral); ok {
		if rightInt, ok := right.(*ast.IntegerLiteral); ok {
			return &ast.BooleanLiteral{Value: leftInt.Value.Cmp(rightInt.Value) <= 0}, nil
		}
	}

	return nil, fmt.Errorf("'<=' cannot be applied to '%s' and '%s'", leftExp, rightExp)
}

func greaterThan(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	left, right, err := evaluateOperands(leftExp, rightExp, scope)
	if err != nil {
		return nil, err
	}

	if leftNum, ok := left.(*ast.NumberLiteral); ok {
		rightVal, err := conformDataType(ast.NUM, right)
		if err != nil {
			return nil, err
		}
		if rightNum, ok := rightVal.(*ast.NumberLiteral); ok {
			return &ast.BooleanLiteral{Value: leftNum.Value > rightNum.Value}, nil
		}
	} else if leftInt, ok := left.(*ast.IntegerLiteral); ok {
		if rightInt, ok := right.(*ast.IntegerLiteral); ok {
			return &ast.BooleanLiteral{Value: leftInt.Value.Cmp(rightInt.Value) > 0}, nil
		}
	}

	return nil, fmt.Errorf("'>' cannot be applied to '%s' and '%s'", leftExp, rightExp)
}

func greaterEqual(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	left, right, err := evaluateOperands(leftExp, rightExp, scope)
	if err != nil {
		return nil, err
	}

	if leftNum, ok := left.(*ast.NumberLiteral); ok {
		rightVal, err := conformDataType(ast.NUM, right)
		if err != nil {
			return nil, err
		}
		if rightNum, ok := rightVal.(*ast.NumberLiteral); ok {
			return &ast.BooleanLiteral{Value: leftNum.Value >= rightNum.Value}, nil
		}
	} else if leftInt, ok := left.(*ast.IntegerLiteral); ok {
		if rightInt, ok := right.(*ast.IntegerLiteral); ok {
			return &ast.BooleanLiteral{Value: leftInt.Value.Cmp(rightInt.Value) >= 0}, nil
		}
	}

	return nil, fmt.Errorf("'>=' cannot be applied to '%s' and '%s'", leftExp, rightExp)
}
