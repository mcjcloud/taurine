package evaluator

import (
	"fmt"

	"github.com/mcjcloud/taurine/ast"
)

// attempts to evaluate an internal function or property (prop) on some type (obj)
func evaluateIntern(obj, prop ast.Expression, scope *Scope) (ast.Expression, error) {
  if strObj, ok := obj.(*ast.StringLiteral); ok {
    return evaluateInternStr(strObj, prop, scope)
  } else if arrObj, ok := obj.(*ast.ArrayExpression); ok {
    return evaluateInternArr(arrObj, prop, scope)
  }
  return nil, fmt.Errorf("'.' cannot be applied to %v", obj)
}

