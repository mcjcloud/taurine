package evaluator

import (
  "fmt"
  "errors"

  "github.com/mcjcloud/taurine/ast"
)

func evaluateInternArr(arr *ast.ArrayExpression, prop ast.Expression, scope *Scope) (ast.Expression, error) {
  if id, ok := prop.(*ast.Identifier); ok {
    switch id.Name {
    case "length":
      return arrLength(arr)
    default:
      return nil, fmt.Errorf("error resolving proprty '%s'", id.Name)
    }
  } else if fn, ok := prop.(*ast.FunctionCall); ok {
    // the function should be an identifier
    if id, ok := fn.Function.(*ast.Identifier); ok {
      switch id.Name {
      case "slice":
        return arrSlice(arr, fn.Arguments, scope)
      case "map":
        return arrMap(arr, fn.Arguments, scope)
      case "forEach":
        return nil, arrForEach(arr, fn.Arguments, scope)
      default:
        return nil, fmt.Errorf("error resolving function '%s'", id.Name)
      }
    }
    return nil, fmt.Errorf("error resolving property '%s'", prop)
  }
  return nil, fmt.Errorf("error resolving property '%s'", prop)
}

// returns the length of the array
func arrLength(arr *ast.ArrayExpression) (*ast.NumberLiteral, error) {
  return &ast.NumberLiteral{
    Value: float64(len(arr.Expressions)),
  }, nil
}

// return a subset range of the array
func arrSlice(arr *ast.ArrayExpression, args []ast.Expression, scope *Scope) (*ast.ArrayExpression, error) {
  if len(args) < 1 || len(args) > 2 {
    return nil, fmt.Errorf("expected 1-2 argument but found %d", len(args))
  }

  var start int
  startExp, err := evaluateExpression(args[0], scope)
  if err != nil {
    return nil, err
  }
  if startNum, ok := startExp.(*ast.NumberLiteral); !ok || startNum.Value != float64(int(startNum.Value)) {
    return nil, fmt.Errorf("expected integer for first argument to slice but found %v", startExp)
  } else {
    start = int(startNum.Value)
  }

  var end int
  if len(args) == 2 {
    endExp, err := evaluateExpression(args[1], scope)
    if err != nil {
      return nil, err
    }
    if endNum, ok := endExp.(*ast.NumberLiteral); !ok || endNum.Value != float64(int(endNum.Value)) {
      return nil, fmt.Errorf("expected integer for second argument to slice but found %v", startExp)
    } else {
      end = int(endNum.Value)
    }
  } else {
    end = len(arr.Expressions)
  }

  // check out of range
  if start < 0 || start > end {
    return nil, fmt.Errorf("start index is outside of range 0-%d", end)
  }
  if end > len(arr.Expressions) {
    return nil, fmt.Errorf("end index is outside of range %d-%d", start, len(arr.Expressions))
  }

  return &ast.ArrayExpression{
    Expressions: arr.Expressions[start:end],
  }, nil
}

// map function for array
func arrMap(arr *ast.ArrayExpression, args []ast.Expression, scope *Scope) (*ast.ArrayExpression, error) {
  if len(args) != 1 {
    return nil, fmt.Errorf("expected 1 argument in map but found %d", len(args))
  }

  fnExp, err := evaluateExpression(args[0], scope)
  if err != nil {
    return nil, err
  }

  if fn, ok := fnExp.(*ScopedFunction); ok {
    // check parameters
    if len(fn.Function.Parameters) == 2 && fn.Function.Parameters[1].SymbolType != ast.NUM {
      return nil, errors.New("expected num for second argument type")
    }
    if len(fn.Function.Parameters) == 3 && fn.Function.Parameters[2].SymbolType != ast.NUM {
      return nil, errors.New("expected num for third argument type")
    }

    // check return type
    if fn.Function.ReturnType == ast.VOID {
      return nil, errors.New("map function must have return type")
    }

    // loop over array expressions
    newArr := make([]ast.Expression, 0)
    for i, exp := range arr.Expressions {
      // build params
      if len(fn.Function.Parameters) > 0 {
        fn.Scope.Set(fn.Function.Parameters[0].Symbol, exp)
      }
      if len(fn.Function.Parameters) > 1 {
        fn.Scope.Set(fn.Function.Parameters[1].Symbol, &ast.NumberLiteral{Value: float64(i)})
      }
      if len(fn.Function.Parameters) > 2 {
        fn.Scope.Set(fn.Function.Parameters[2].Symbol, &ast.NumberLiteral{Value: float64(len(arr.Expressions))})
      }

      // call function
      if err := executeStatement(fn.Function.Body, fn.Scope); err != nil {
        return nil, err
      }
      newArr = append(newArr, fn.Scope.ReturnValue)
    }

    // return the resultant array
    return &ast.ArrayExpression{
      Expressions: newArr,
    }, nil
  }
  return nil, fmt.Errorf("expected function argument to map")
}

func arrForEach(arr *ast.ArrayExpression, args []ast.Expression, scope *Scope) error {
  _, err := arrMap(arr, args, scope)
  return err
}

