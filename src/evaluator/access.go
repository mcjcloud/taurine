package evaluator

import (
  "errors"
  "fmt"

  "github.com/mcjcloud/taurine/ast"
)

func arrayIndex(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
  left, right, err := evaluateOperands(leftExp, rightExp, scope)
  if err != nil {
    return nil, err
  }

  if leftArr, ok := left.(*ast.ArrayExpression); ok {
    if rightNum, ok := right.(*ast.NumberLiteral); ok {
      if i := int(rightNum.Value); float64(i) == rightNum.Value {
        if i < 0 || i > len(leftArr.Expressions) {
          return nil, fmt.Errorf("index %d out of range", i)
        }
        return evaluateExpression(leftArr.Expressions[i], scope)
      } else {
        return nil, errors.New("'@' index must evalute to an integer")
      }
    }
  } else if leftStr, ok := left.(*ast.StringLiteral); ok {
    if rightNum, ok := right.(*ast.NumberLiteral); ok {
      if i := int(rightNum.Value); float64(i) == rightNum.Value {
        if i < 0 || i > len(leftStr.Value) {
          return nil, fmt.Errorf("index %d out of range", i)
        }
        return &ast.StringLiteral{Value: string([]rune(leftStr.Value)[i])}, nil
      } else {
        return nil, errors.New("'@' index must evaluate to an integer")
      }
    }
  }
  return nil, errors.New("'@' operator must be in form arr@integer")
}

func createRange(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
  left, right, err := evaluateOperands(leftExp, rightExp, scope)
  if err != nil {
    return nil, err
  }

  // make sure each left and right operator are integers
  if leftNum, ok := left.(*ast.NumberLiteral); ok && leftNum.Value == float64(int(leftNum.Value)) {
    if rightNum, ok := right.(*ast.NumberLiteral); ok && rightNum.Value == float64(int(rightNum.Value)) {
      var direction int
      if leftNum.Value < rightNum.Value {
        direction = 1
      } else if rightNum.Value > leftNum.Value {
        direction = -1
      } else {
        return &ast.ArrayExpression{
          Expressions: []ast.Expression{leftNum},
        }, nil
      }
      // use direction to iterate and populate array
      arr := make([]ast.Expression, 0)
      for i := int(leftNum.Value); i != int(rightNum.Value); i += direction {
        arr = append(arr, &ast.NumberLiteral{Value: float64(i)})
      }
      return &ast.ArrayExpression{Expressions: arr}, nil
    }
  }
  return nil, errors.New("'..' must have operands of type integer")
}

func dot(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
  left, err := evaluateExpression(leftExp, scope)
  if err != nil {
    return nil, fmt.Errorf("error accessing obj member: %s", err.Error())
  }
  if leftObj, ok := left.(*ast.ObjectLiteral); ok {
    // the right side must be either an identifier, fn call, or another dot operator
    if rightIdentifier, ok := rightExp.(*ast.Identifier); ok {
      if leftObj.Value[rightIdentifier.Name] != nil {
        return evaluateExpression(leftObj.Value[rightIdentifier.Name], scope)
      }
    } else if rightFnCall, ok := rightExp.(*ast.FunctionCall); ok {
      objScope := NewScopeOfObject(leftObj, scope)
      return evaluateFunctionCall(rightFnCall, objScope)
    } else if rightDotOp, ok := rightExp.(*ast.OperationExpression); ok && rightDotOp.Operator == ast.DOT {
      // create a new scope with the parent obj as scope
      objScope := NewScopeOfObject(leftObj, scope)
      return evaluateOperation(rightDotOp, objScope)
    } else if rightAsgn, ok := rightExp.(*ast.AssignmentExpression); ok {
      // TODO: consider moving this logic. Perhaps the assignment should be moved
      // further up the tree and this logic should be handled by the assignment fn
      // evaluate the new value and update the object property
      newVal, err := evaluateExpression(rightAsgn.Value, scope)
      if err != nil {
        return nil, err
      }
      leftObj.Value[rightAsgn.Identifier.Name] = newVal
      return newVal, nil
    }
    return nil, errors.New("right side of '.' must be identifier or function call")
  } else {
    return evaluateIntern(left, rightExp, scope)
  }
}

