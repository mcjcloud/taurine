package evaluator

import (
	"errors"
	"fmt"

	"github.com/mcjcloud/taurine/ast"
)

func add(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
  left, right, err := evaluateOperands(leftExp, rightExp, scope)
  if err != nil {
    return nil, err
  }

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

	if leftNum, ok := left.(*ast.NumberLiteral); ok {
		if rightNum, ok := right.(*ast.NumberLiteral); ok {
      if rightNum.Value == float64(0) {
        return nil, errors.New("divide by 0 error")
      }
			return &ast.NumberLiteral{Value: leftNum.Value / rightNum.Value}, nil
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

	if leftNum, ok := left.(*ast.NumberLiteral); ok && leftNum.Value == float64(int(leftNum.Value)) {
		if rightNum, ok := right.(*ast.NumberLiteral); ok && rightNum.Value == float64(int(rightNum.Value)) {
      if rightNum.Value == float64(0) {
        return nil, errors.New("divide by 0 error")
      }
			return &ast.NumberLiteral{Value: float64(int(leftNum.Value) % int(rightNum.Value))}, nil
		}
	}
	return nil, errors.New("'%' operator only applies to type num")
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

