package evaluator

import (
  "fmt"

  "github.com/mcjcloud/taurine/ast"
)

func equalEqual(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
  left, right, err := evaluateOperands(leftExp, rightExp, scope)
  if err != nil {
    return nil, err
  }

  leftNum, lok := left.(*ast.NumberLiteral)
  rightNum, rok := right.(*ast.NumberLiteral)
  if lok && rok {
    return &ast.BooleanLiteral{Value: leftNum.Value == rightNum.Value}, nil
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
  rightNum, rok := right.(*ast.NumberLiteral)
  if lok && rok {
    return &ast.BooleanLiteral{Value: leftNum.Value != rightNum.Value}, nil
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
    if rightNum, ok := right.(*ast.NumberLiteral); ok {
      return &ast.BooleanLiteral{Value: leftNum.Value < rightNum.Value}, nil
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
    if rightNum, ok := right.(*ast.NumberLiteral); ok {
      return &ast.BooleanLiteral{Value: leftNum.Value <= rightNum.Value}, nil
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
    if rightNum, ok := right.(*ast.NumberLiteral); ok {
      return &ast.BooleanLiteral{Value: leftNum.Value > rightNum.Value}, nil
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
    if rightNum, ok := right.(*ast.NumberLiteral); ok {
      return &ast.BooleanLiteral{Value: leftNum.Value >= rightNum.Value}, nil
    }
  }
  return nil, fmt.Errorf("'>=' cannot be applied to '%s' and '%s'", leftExp, rightExp)
}

