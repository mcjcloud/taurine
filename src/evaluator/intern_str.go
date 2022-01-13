package evaluator

import (
	"fmt"
  "strings"

	"github.com/mcjcloud/taurine/ast"
)

func evaluateInternStr(str *ast.StringLiteral, prop ast.Expression, scope *Scope) (ast.Expression, error) {
  if id, ok := prop.(*ast.Identifier); ok {
    switch id.Name {
    case "length":
      return strLength(str)
    default:
      return nil, fmt.Errorf("error resolving property '%s'", id.Name)
    }
  } else if fn, ok := prop.(*ast.FunctionCall); ok {
    // the function should be an identifier
    if id, ok := fn.Function.(*ast.Identifier); ok {
      switch id.Name {
      case "toUpperCase":
        return strToUpperCase(str)
      case "toLowerCase":
        return strToLowerCase(str)
      case "toArray":
        return strToArray(str)
      case "substr":
        return strSubstr(str, fn.Arguments, scope)
      default:
        return nil, fmt.Errorf("error resolving property '%s'", id.Name)
      }
    }
    return nil, fmt.Errorf("error resolving property '%s'", prop)
  }
  return nil, fmt.Errorf("error resolving property '%s'", prop)
}

// access the length of the string
func strLength(str *ast.StringLiteral) (*ast.NumberLiteral, error) {
  return &ast.NumberLiteral{
    Value: float64(len(str.Value)),
  }, nil
}

// convert the string to uppercase
func strToUpperCase(str *ast.StringLiteral) (*ast.StringLiteral, error) {
  return &ast.StringLiteral{
    Value: strings.ToUpper(str.Value),
  }, nil
}

// convert the string to lowercase
func strToLowerCase(str *ast.StringLiteral) (*ast.StringLiteral, error) {
  return &ast.StringLiteral{
    Value: strings.ToLower(str.Value),
  }, nil
}

// convert the string to an array of its characters
func strToArray(str *ast.StringLiteral) (*ast.ArrayExpression, error) {
  res := make([]ast.Expression, 0)
  for _, c := range str.Value {
    res = append(res, &ast.StringLiteral{Value: string(c)})
  }
  return &ast.ArrayExpression{
    Expressions: res,
  }, nil
}

// return a substring of the given string
func strSubstr(str *ast.StringLiteral, args []ast.Expression, scope *Scope) (*ast.StringLiteral, error) {
  if len(args) < 1 || len(args) > 2 {
    return nil, fmt.Errorf("expected 1-2 argument but found %d", len(args))
  }

  var start int
  startExp, err := evaluateExpression(args[0], scope)
  if err != nil {
    return nil, err
  }
  if startNum, ok := startExp.(*ast.NumberLiteral); !ok || startNum.Value != float64(int(startNum.Value)) {
    return nil, fmt.Errorf("expected integer for first argument to substr but found %v", startExp)
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
      return nil, fmt.Errorf("expected integer for second argument to substr but found %v", startExp)
    } else {
      end = int(endNum.Value)
    }
  } else {
    end = len(str.Value)
  }

  // check out of range
  if start < 0 || start > end {
    return nil, fmt.Errorf("start index is outside of range 0-%d", end)
  }
  if end > len(str.Value) {
    return nil, fmt.Errorf("end index is outside of range %d-%d", start, len(str.Value))
  }

  return &ast.StringLiteral{
    Value: str.Value[start:end],
  }, nil
}

