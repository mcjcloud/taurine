package evaluator

import (
  "errors"
  "fmt"
  "math/big"

  "github.com/mcjcloud/taurine/ast"
)

func arrayIndex(leftExp, rightExp ast.Expression, scope *Scope) (ast.Expression, error) {
  left, right, err := evaluateOperands(leftExp, rightExp, scope)
  if err != nil {
    return nil, err
  }

  if leftArr, ok := left.(*ast.ArrayExpression); ok {
    if rightNum, ok := right.(*ast.IntegerLiteral); ok {
      i := int(rightNum.Value.Int64())
      if i < 0 || i > len(leftArr.Expressions) {
        return nil, fmt.Errorf("index %d out of range", i)
      }
      return evaluateExpression(leftArr.Expressions[i], scope)
    }
  } else if leftStr, ok := left.(*ast.StringLiteral); ok {
    if rightNum, ok := right.(*ast.IntegerLiteral); ok {
      i := int(rightNum.Value.Int64())
      if i < 0 || i > len(leftStr.Value) {
        return nil, fmt.Errorf("index %d out of range", i)
      }
      return &ast.StringLiteral{Value: string([]rune(leftStr.Value)[i])}, nil
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
  if leftNum, ok := left.(*ast.IntegerLiteral); ok {
    if rightNum, ok := right.(*ast.IntegerLiteral); ok {
      var direction int
      if leftNum.Value.Cmp(rightNum.Value) < 0 {
        direction = 1
      } else if leftNum.Value.Cmp(rightNum.Value) > 0 {
        direction = -1
      } else {
        return &ast.ArrayExpression{
          Expressions: []ast.Expression{leftNum},
        }, nil
      }
      // use direction to iterate and populate array
      arr := make([]ast.Expression, 0)
      for i := int(leftNum.Value.Int64()); i != int(rightNum.Value.Int64()); i += direction {
        arr = append(arr, &ast.IntegerLiteral{Value: big.NewInt(int64(i))})
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

