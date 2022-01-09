package evaluator

import (
  "errors"
  "fmt"

  "github.com/mcjcloud/taurine/ast"
)

func evaluateOperation(op *ast.OperationExpression, scope *Scope) (ast.Expression, error) {
  // check . operator first because it doesn't involve evaluating the right side immediately
  if op.Operator == ast.DOT { // "."
    left, err := evaluateExpression(op.LeftExpression, scope)
    if err != nil {
      return nil, fmt.Errorf("error accessing obj member: %s", err.Error())
    }
    if leftObj, ok := left.(*ast.ObjectLiteral); ok {
      // the right side must be either an identifier, fn call, or another dot operator
      if rightIdentifier, ok := op.RightExpression.(*ast.Identifier); ok {
        if leftObj.Value[rightIdentifier.Name] != nil {
          return evaluateExpression(leftObj.Value[rightIdentifier.Name], scope)
        }
      } else if rightFnCall, ok := op.RightExpression.(*ast.FunctionCall); ok {
        objScope := NewScopeOfObject(leftObj, scope)
        return evaluateFunctionCall(rightFnCall, objScope)
      } else if rightDotOp, ok := op.RightExpression.(*ast.OperationExpression); ok && rightDotOp.Operator == ast.DOT {
        // create a new scope with the parent obj as scope
        objScope := NewScopeOfObject(leftObj, scope)
        return evaluateOperation(rightDotOp, objScope)
      }
      return nil, errors.New("right side of '.' must be identifier or function call")
    }
    return nil, errors.New("'.' operator must be in form obj.identifier")
  }

  left, err := evaluateExpression(op.LeftExpression, scope)
  if err != nil {
    return nil, err
  }
  right, err := evaluateExpression(op.RightExpression, scope)
  if err != nil {
    return nil, err
  }
  if op.Operator == ast.PLUS { // "+"
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
  } else if op.Operator == ast.MINUS { // "-"
    if leftNum, ok := left.(*ast.NumberLiteral); ok {
      if rightNum, ok := right.(*ast.NumberLiteral); ok {
        return &ast.NumberLiteral{Value: leftNum.Value - rightNum.Value}, nil
      }
    }
    return nil, errors.New("'-' operator only applies to type num")
  } else if op.Operator == ast.DIVIDE { // "/"
    if leftNum, ok := left.(*ast.NumberLiteral); ok {
      if rightNum, ok := right.(*ast.NumberLiteral); ok {
        return &ast.NumberLiteral{Value: leftNum.Value / rightNum.Value}, nil
      }
    }
    return nil, errors.New("'/' operator only applies to type num")
  } else if op.Operator == ast.MULTIPLY { // "*"
    if leftNum, ok := left.(*ast.NumberLiteral); ok {
      if rightNum, ok := right.(*ast.NumberLiteral); ok {
        return &ast.NumberLiteral{Value: leftNum.Value * rightNum.Value}, nil
      }
    }
    return nil, errors.New("'*' operator only applies to type num")
  } else if op.Operator == "==" {
    leftNum, lok := left.(*ast.NumberLiteral)
    rightNum, rok := right.(*ast.NumberLiteral)
    if (lok && !rok) || (!lok && rok) {
      return nil, errors.New("'==' operator cannot be applied to arguments of different types")
    } else if lok && rok {
      return &ast.BooleanLiteral{Value: leftNum.Value == rightNum.Value}, nil
    }

    leftStr, lok := left.(*ast.StringLiteral)
    rightStr, rok := right.(*ast.StringLiteral)
    if (lok && !rok) || (!lok && rok) {
      return nil, errors.New("'==' operator cannot be applied to arguments of different types")
    } else if lok && rok {
      return &ast.BooleanLiteral{Value: leftStr.Value == rightStr.Value}, nil
    }

    leftBool, lok := left.(*ast.BooleanLiteral)
    rightBool, rok := right.(*ast.BooleanLiteral)
    if (lok && !rok) || (!lok && rok) {
      return nil, errors.New("'==' operator cannot be applied to arguments of different types")
    } else if lok && rok {
      return &ast.BooleanLiteral{Value: leftBool.Value == rightBool.Value}, nil
    }
  } else if op.Operator == "!=" {
    leftNum, lok := left.(*ast.NumberLiteral)
    rightNum, rok := right.(*ast.NumberLiteral)
    if (lok && !rok) || (!lok && rok) {
      return nil, errors.New("'!=' operator cannot be applied to arguments of different types")
    } else if lok && rok {
      return &ast.BooleanLiteral{Value: leftNum.Value != rightNum.Value}, nil
    }

    leftStr, lok := left.(*ast.StringLiteral)
    rightStr, rok := right.(*ast.StringLiteral)
    if (lok && !rok) || (!lok && rok) {
      return nil, errors.New("'!=' operator cannot be applied to arguments of different types")
    } else if lok && rok {
      return &ast.BooleanLiteral{Value: leftStr.Value != rightStr.Value}, nil
    }

    leftBool, lok := left.(*ast.BooleanLiteral)
    rightBool, rok := right.(*ast.BooleanLiteral)
    if (lok && !rok) || (!lok && rok) {
      return nil, errors.New("'!=' operator cannot be applied to arguments of different types")
    } else if lok && rok {
      return &ast.BooleanLiteral{Value: leftBool.Value != rightBool.Value}, nil
    }
  } else if op.Operator == "<" {
    if leftNum, ok := left.(*ast.NumberLiteral); ok {
      if rightNum, ok := right.(*ast.NumberLiteral); ok {
        return &ast.BooleanLiteral{Value: leftNum.Value < rightNum.Value}, nil
      }
    }
    return nil, errors.New("'<' operator only applies to type num")
  } else if op.Operator == ">" {
    if leftNum, ok := left.(*ast.NumberLiteral); ok {
      if rightNum, ok := right.(*ast.NumberLiteral); ok {
        return &ast.BooleanLiteral{Value: leftNum.Value > rightNum.Value}, nil
      }
    }
    return nil, errors.New("'>' operator only applies to type num")
  } else if op.Operator == "<=" {
    if leftNum, ok := left.(*ast.NumberLiteral); ok {
      if rightNum, ok := right.(*ast.NumberLiteral); ok {
        return &ast.BooleanLiteral{Value: leftNum.Value <= rightNum.Value}, nil
      }
    }
    return nil, errors.New("'<=' operator only applies to type num")
  } else if op.Operator == ">=" {
    if leftNum, ok := left.(*ast.NumberLiteral); ok {
      if rightNum, ok := right.(*ast.NumberLiteral); ok {
        return &ast.BooleanLiteral{Value: leftNum.Value >= rightNum.Value}, nil
      }
    }
    return nil, errors.New("'>=' operator only applies to type num")
  } else if op.Operator == "@" {
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
  } else if op.Operator == ast.RANGE { // ".."
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
  return nil, fmt.Errorf("unrecognized operator '%s'", op.Operator)
}

func builtInLen(exp ast.Expression, scope *Scope) (*ast.NumberLiteral, error) {
  evExp, err := evaluateExpression(exp, scope)
  if err != nil {
    return nil, err
  }
  if strExp, ok := evExp.(*ast.StringLiteral); ok {
    return &ast.NumberLiteral{Value: float64(len(strExp.Value))}, nil
  } else if arrExp, ok := evExp.(*ast.ArrayExpression); ok {
    return &ast.NumberLiteral{Value: float64(len(arrExp.Expressions))}, nil
  }
  return nil, errors.New("len can only be called on type str or arr")
}

