package evaluator

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/mcjcloud/taurine/pkg/ast"
)

func add(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	left, right, err := evaluateOperands(leftExp, rightExp, scope)
	if err != nil {
		return nil, err
	}

	if leftNum, ok := left.(*ast.NumberLiteral); ok {

		// add either num, int, or string
		if rightNum, ok := right.(*ast.NumberLiteral); ok {
			return &ast.NumberLiteral{Value: leftNum.Value + rightNum.Value}, nil
		} else if rightInt, ok := right.(*ast.IntegerLiteral); ok {
			return &ast.NumberLiteral{Value: leftNum.Value + float64(rightInt.Value.Int64())}, nil
		} else if rightStr, ok := right.(*ast.StringLiteral); ok {
			return &ast.StringLiteral{Value: fmt.Sprintf("%f%s", leftNum.Value, rightStr.Value)}, nil
		}
	} else if leftInt, ok := left.(*ast.IntegerLiteral); ok {

		// add either num, int, or string
		if rightNum, ok := right.(*ast.NumberLiteral); ok {
			return &ast.NumberLiteral{Value: float64(leftInt.Value.Int64()) + rightNum.Value}, nil
		} else if rightInt, ok := right.(*ast.IntegerLiteral); ok {
			newInt := new(big.Int).Add(leftInt.Value, rightInt.Value)
			return &ast.IntegerLiteral{Value: newInt}, nil
		} else if rightStr, ok := right.(*ast.StringLiteral); ok {
			return &ast.StringLiteral{Value: fmt.Sprintf("%s%s", leftInt.Value, rightStr.Value)}, nil
		}
	} else if leftStr, ok := left.(*ast.StringLiteral); ok {

		// add stringified version of whatever is on right side
		return &ast.StringLiteral{Value: fmt.Sprintf("%s%s", leftStr.String(), right.String())}, nil
	}
	return nil, fmt.Errorf("'+' operator is not applicable to arguments %s and %s", leftExp, rightExp)
}

func addAndAssign(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	leftId, err := assertIdentifier(leftExp)
	if err != nil {
		return nil, err
	}

	res, err := add(leftId, rightExp, scope)
	if err != nil {
		return nil, err
	}
	scope.Set(leftId.Name, res)

	return res, nil
}

func minus(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	left, right, err := evaluateOperands(leftExp, rightExp, scope)
	if err != nil {
		return nil, err
	}

	if leftNum, ok := left.(*ast.NumberLiteral); ok {
		if rightNum, ok := right.(*ast.NumberLiteral); ok {
			return &ast.NumberLiteral{Value: leftNum.Value - rightNum.Value}, nil
		} else if rightInt, ok := right.(*ast.IntegerLiteral); ok {
			return &ast.NumberLiteral{Value: leftNum.Value - float64(rightInt.Value.Int64())}, nil
		}
	} else if leftInt, ok := left.(*ast.IntegerLiteral); ok {
		if rightNum, ok := right.(*ast.NumberLiteral); ok {
			return &ast.NumberLiteral{Value: float64(leftInt.Value.Int64()) - rightNum.Value}, nil
		} else if rightInt, ok := right.(*ast.IntegerLiteral); ok {
			newInt := new(big.Int).Sub(leftInt.Value, rightInt.Value)
			return &ast.IntegerLiteral{Value: newInt}, nil
		}
	}

	return nil, errors.New("'-' operator only applies to type num")
}

func minusAndAssign(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	leftId, err := assertIdentifier(leftExp)
	if err != nil {
		return nil, err
	}

	res, err := minus(leftId, rightExp, scope)
	if err != nil {
		return nil, err
	}
	scope.Set(leftId.Name, res)

	return res, nil
}

func multiply(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	left, right, err := evaluateOperands(leftExp, rightExp, scope)
	if err != nil {
		return nil, err
	}

	if leftNum, ok := left.(*ast.NumberLiteral); ok {
		if rightNum, ok := right.(*ast.NumberLiteral); ok {
			return &ast.NumberLiteral{Value: leftNum.Value * rightNum.Value}, nil
		} else if rightInt, ok := right.(*ast.IntegerLiteral); ok {
			return &ast.NumberLiteral{Value: leftNum.Value * float64(rightInt.Value.Int64())}, nil
		}
	} else if leftInt, ok := left.(*ast.IntegerLiteral); ok {
		if rightNum, ok := right.(*ast.NumberLiteral); ok {
			return &ast.NumberLiteral{Value: float64(leftInt.Value.Int64()) * rightNum.Value}, nil
		} else if rightInt, ok := right.(*ast.IntegerLiteral); ok {
			newInt := new(big.Int).Mul(leftInt.Value, rightInt.Value)
			return &ast.IntegerLiteral{Value: newInt}, nil
		}
	}

	return nil, errors.New("'*' operator only applies to type num")
}

func multiplyAndAssign(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	leftId, err := assertIdentifier(leftExp)
	if err != nil {
		return nil, err
	}

	res, err := multiply(leftId, rightExp, scope)
	if err != nil {
		return nil, err
	}
	scope.Set(leftId.Name, res)

	return res, nil
}

func divide(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	left, right, err := evaluateOperands(leftExp, rightExp, scope)
	if err != nil {
		return nil, err
	}

	if rightNum, ok := right.(*ast.NumberLiteral); ok {
		if rightNum.Value == float64(0) {
			return nil, errors.New("divide by 0 error")
		}
		if leftNum, ok := left.(*ast.NumberLiteral); ok {
			return &ast.NumberLiteral{Value: leftNum.Value / rightNum.Value}, nil
		} else if leftInt, ok := left.(*ast.IntegerLiteral); ok {
			return &ast.NumberLiteral{Value: float64(leftInt.Value.Int64()) / rightNum.Value}, nil
		}
	} else if rightInt, ok := right.(*ast.IntegerLiteral); ok {
		if rightInt.Value.Int64() == 0 {
			return nil, errors.New("divide by 0 error")
		}
		if leftNum, ok := left.(*ast.NumberLiteral); ok {
			return &ast.NumberLiteral{Value: leftNum.Value / float64(rightInt.Value.Int64())}, nil
		} else if leftInt, ok := left.(*ast.IntegerLiteral); ok {
			newInt := new(big.Int).Div(leftInt.Value, rightInt.Value)
			return &ast.IntegerLiteral{Value: newInt}, nil
		}
	}

	return nil, errors.New("'/' operator only applies to type num")
}

func divideAndAssign(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	leftId, err := assertIdentifier(leftExp)
	if err != nil {
		return nil, err
	}

	res, err := divide(leftId, rightExp, scope)
	if err != nil {
		return nil, err
	}
	scope.Set(leftId.Name, res)

	return res, nil
}

func modulo(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	left, right, err := evaluateOperands(leftExp, rightExp, scope)
	if err != nil {
		return nil, err
	}

	if leftInt, ok := left.(*ast.IntegerLiteral); ok {
		if rightInt, ok := right.(*ast.IntegerLiteral); ok {
			if rightInt.Value.Int64() == 0 {
				return nil, errors.New("divide by 0 error")
			}
			newInt := new(big.Int).Mod(leftInt.Value, rightInt.Value)
			return &ast.IntegerLiteral{Value: newInt}, nil
		}
	}
	return nil, errors.New("'%' operator only applies to integers")
}

func moduloAndAssign(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
	leftId, err := assertIdentifier(leftExp)
	if err != nil {
		return nil, err
	}

	res, err := modulo(leftId, rightExp, scope)
	if err != nil {
		return nil, err
	}
	scope.Set(leftId.Name, res)

	return res, nil
}
