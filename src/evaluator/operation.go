package evaluator

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/mcjcloud/taurine/ast"
)

func assertIdentifier(exp ast.Expression) (*ast.Identifier, error) {
  if id, ok := exp.(*ast.Identifier); ok {
    return id, nil
  }
  return nil, fmt.Errorf("expected identifier but found %v", exp)
}

func evaluateOperands(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, ast.Expression, error) {
  left, err := evaluateExpression(leftExp, scope)
	if err != nil {
		return nil, nil, err
	}
	right, err := evaluateExpression(rightExp, scope)
	if err != nil {
		return nil, nil, err
	}
  return left, right, nil
}

func evaluateOperation(op *ast.OperationExpression, scope *Scope) (ast.Expression, error) {
	left := op.LeftExpression
	right := op.RightExpression
  switch op.Operator {
  case ast.PLUS:
    return add(left, right, scope)
  case ast.PLUS_EQUAL:
    return addAndAssign(left, right, scope)
  case ast.MINUS:
    return minus(left, right, scope)
  case ast.MINUS_EQUAL:
    return minusAndAssign(left, right, scope)
  case ast.MULTIPLY:
    return multiply(left, right, scope)
  case ast.MULTIPLY_EQUAL:
    return multiplyAndAssign(left, right, scope)
  case ast.DIVIDE:
    return divide(left, right, scope)
  case ast.DIVIDE_EQUAL:
    return divideAndAssign(left, right, scope)
  case ast.MODULO:
    return modulo(left, right, scope)
  case ast.MODULO_EQUAL:
    return moduloAndAssign(left, right, scope)
  case ast.EQUAL_EQUAL:
    return equalEqual(left, right, scope)
  case ast.NOT_EQUAL:
    return notEqual(left, right, scope)
  case ast.LESS_THAN:
    return lessThan(left, right, scope)
  case ast.LESS_EQUAL:
    return lessEqual(left, right, scope)
  case ast.GREATER_THAN:
    return greaterThan(left, right, scope)
  case ast.GREATER_EQUAL:
    return greaterEqual(left, right, scope)
  case ast.AT:
    return arrayIndex(left, right, scope)
  case ast.RANGE:
    return createRange(left, right, scope)
  case ast.DOT:
    return dot(left, right, scope)
  default:
    return nil, fmt.Errorf("unrecognized operator '%s'", op.Operator)
  }
}

func builtInLen(exp ast.Expression, scope *Scope) (*ast.IntegerLiteral, error) {
	evExp, err := evaluateExpression(exp, scope)
	if err != nil {
		return nil, err
	}
	if strExp, ok := evExp.(*ast.StringLiteral); ok {
		return &ast.IntegerLiteral{Value: big.NewInt(int64(len(strExp.Value)))}, nil
	} else if arrExp, ok := evExp.(*ast.ArrayExpression); ok {
		return &ast.IntegerLiteral{Value: big.NewInt(int64(len(arrExp.Expressions)))}, nil
	}
	return nil, errors.New("len can only be called on type str or arr")
}

func builtInInt(exp ast.Expression, scope *Scope) (*ast.NumberLiteral, error) {
  if expEv, err := evaluateExpression(exp, scope); err == nil {
    if num, ok := expEv.(*ast.NumberLiteral); ok {
      return &ast.NumberLiteral{Value: float64(int(num.Value))}, nil
    }
  }
  return nil, errors.New("int() can only be called on type num")
}

