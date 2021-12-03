package evaluator

import (
  "bufio"
  "errors"
  "fmt"
  "os"
  "strings"

  "github.com/mcjcloud/taurine/ast"
)

// Evaluate evaluates the code and does stuff
func Evaluate(block *ast.BlockStatement) error {
  scope := NewScope()
  for _, stmt := range block.Statements {
    // do the statement
    err := executeStatement(stmt, scope)
    if err != nil {
      return err
    }
  }
  return nil
}

func executeStatement(stmt ast.Statement, scope *Scope) error {
  if etchStmt, ok := stmt.(*ast.EtchStatement); ok {
    if err := executeEtchStatement(etchStmt, scope); err != nil {
      return err
    }
  } else if readStmt, ok := stmt.(*ast.ReadStatement); ok {
    if err := executeReadStatement(readStmt, scope); err != nil {
      return err
    }
  } else if declStmt, ok := stmt.(*ast.VariableDecleration); ok {
    if err := executeVariableDecleration(declStmt, scope); err != nil {
      return err
    }
  } else if funcStmt, ok := stmt.(*ast.FunctionDecleration); ok {
    if err := executeFunctionDecleration(funcStmt, scope); err != nil {
      return err
    }
  } else if expStmt, ok := stmt.(*ast.ExpressionStatement); ok {
    _, err := evaluateExpression(expStmt.Expression, scope)
    return err
  } else if blockStmt, ok := stmt.(*ast.BlockStatement); ok {
    subScope := NewScopeWithParent(scope)
    for _, s := range blockStmt.Statements {
      err := executeStatement(s, subScope)
      if err != nil {
        return err
      }
      // if a block exists within the current scope, the return value should propogate up
      if subScope.ReturnValue != nil {
        scope.ReturnValue = subScope.ReturnValue
        break
      }
    }
  } else if ifStmt, ok := stmt.(*ast.IfStatement); ok {
    exp, err := evaluateExpression(ifStmt.Condition, scope)
    if err != nil {
      return err
    }
    if boolExp, ok := exp.(*ast.BooleanLiteral); ok {
      if boolExp.Value {
        if err := executeStatement(ifStmt.Statement, scope); err != nil {
          return err
        }
      }
      return nil
    }
    return errors.New("if expression must evaluate to boolean")
  } else if whileStmt, ok := stmt.(*ast.WhileLoopStatement); ok {
    exp, err := evaluateExpression(whileStmt.Condition, scope)
    if err != nil {
      return err
    }
    if boolExp, ok := exp.(*ast.BooleanLiteral); ok {
      subScope := NewScopeWithParent(scope)
      for boolExp.Value {
        err := executeStatement(whileStmt.Statement, subScope)
        if err != nil {
          return err
        }
        exp, err = evaluateExpression(whileStmt.Condition, subScope)
        if err != nil {
          return err
        }
        boolExp, ok = exp.(*ast.BooleanLiteral)
        if !ok {
          return errors.New("while expression is no longer boolean")
        }
        // if there is a return value, the loop should end
        if subScope.ReturnValue != nil {
          scope.ReturnValue = subScope.ReturnValue
          break
        }
      }
    } else {
      return errors.New("while expression must evaluate to boolean")
    }
  } else if rtnStmt, ok := stmt.(*ast.ReturnStatement); ok {
    exp, err := evaluateExpression(rtnStmt.Value, scope)
    if err != nil {
      return err
    }
    scope.ReturnValue = exp
    return nil
  } else {
    return errors.New("unrecognized statement")
  }
  return nil
}

func executeEtchStatement(stmt *ast.EtchStatement, scope *Scope) error {
  var toEtch []string
  for _, exp := range stmt.Expressions {
    if numExp, ok := exp.(*ast.NumberLiteral); ok {
      toEtch = append(toEtch, numExp.String())
    } else if strExp, ok := exp.(*ast.StringLiteral); ok {
      toEtch = append(toEtch, strExp.String())
    } else if idExp, ok := exp.(*ast.Identifier); ok {
      idVal := scope.Get(idExp.Name)
      if idVal != nil {
        toEtch = append(toEtch, idVal.String())
      } else {
        toEtch = append(toEtch, "nil")
      }
    } else {
      expEval, err := evaluateExpression(exp, scope)
      if err != nil {
        return err
      }
      toEtch = append(toEtch, expEval.String())
    }
  }
  fmt.Println(strings.Join(toEtch, " "))
  return nil
}

func executeReadStatement(stmt *ast.ReadStatement, scope *Scope) error {
  if stmt.Prompt != nil {
    fmt.Printf("%s", stmt.Prompt)
  }
  scanner := bufio.NewScanner(os.Stdin)
  if scanner.Scan() {
    scope.Set(stmt.Identifier.Name, &ast.StringLiteral{Value: scanner.Text()})
  } else {
    return errors.New("error reading input")
  }
  return nil
}

func executeVariableDecleration(stmt *ast.VariableDecleration, scope *Scope) error {
  val, err := evaluateExpression(stmt.Value, scope)
  if err != nil {
    return err
  }
  if scope.Variables[stmt.Symbol] != nil {
    return fmt.Errorf("variable '%s' already exists", stmt.Symbol)
  }
  scope.Set(stmt.Symbol, val)
  return nil
}

func executeFunctionDecleration(stmt *ast.FunctionDecleration, scope *Scope) error {
  if scope.Functions[stmt.Symbol] != nil {
    return fmt.Errorf("function '%s' already exists", stmt.Symbol)
  }
  scope.SetFunction(stmt.Symbol, stmt)
  return nil
}

func evaluateExpression(exp ast.Expression, scope *Scope) (ast.Expression, error) {
  if op, ok := exp.(*ast.OperationExpression); ok {
    return evaluateOperation(op, scope)
  } else if id, ok := exp.(*ast.Identifier); ok {
    return scope.Get(id.Name), nil
  } else if asn, ok := exp.(*ast.AssignmentExpression); ok {
    // make sure the identifier exists
    if scope.Get(asn.Identifier.Name) == nil {
      return nil, fmt.Errorf("'%s' was not declared", asn.Identifier.Name)
    }
    val, err := evaluateExpression(asn.Value, scope)
    if err != nil {
      return nil, err
    }
    // update the scope and return the evaluated value
    scope.Set(asn.Identifier.Name, val)
    return val, nil
  } else if fnCall, ok := exp.(*ast.FunctionCall); ok {
    return evaluateFunctionCall(fnCall, scope)
  } else if grpExp, ok := exp.(*ast.GroupExpression); ok {
    return evaluateExpression(grpExp.Expression, scope)
  } else if arrExp, ok := exp.(*ast.ArrayExpression); ok {
    return evaluateArrayExpression(arrExp, scope)
  }
  return exp, nil
}

func evaluateArrayExpression(arr *ast.ArrayExpression, scope *Scope) (ast.Expression, error) {
  exp := &ast.ArrayExpression{Expressions: make([]ast.Expression, len(arr.Expressions))}
  for i, el := range arr.Expressions {
    val, err := evaluateExpression(el, scope)
    if err != nil {
      return nil, err
    }
    exp.Expressions[i] = val
  }
  return exp, nil
}

func evaluateFunctionCall(call *ast.FunctionCall, scope *Scope) (ast.Expression, error) {
  fn := scope.GetFunction(call.Function)
  if fn == nil {
    // TODO: make this cleaner, maybe move built-in functions someplace else
    if call.Function == "len" {
      if len(call.Arguments) != 1 {
        return nil, errors.New("len takes only one argument")
      }
      return builtInLen(call.Arguments[0], scope)
    }
    return nil, fmt.Errorf("function name '%s' was never declared", call.Function)
  }
  decl := fn.Decleration
  if len(decl.Parameters) != len(call.Arguments) {
    return nil, fmt.Errorf("expected '%d' arguments but got '%d' for call to '%s'", len(decl.Parameters), len(call.Arguments), call.Function)
  }
  // evaluate arguments and populate scope
  for i, arg := range call.Arguments {
    exp, err := evaluateExpression(arg, scope)
    if err != nil {
      return nil, err
    }
    // TODO: create a good way to compare data type of argument of parameter
    fn.Scope.Set(decl.Parameters[i].Symbol, exp)
  }
  // execute statements
  if err := executeStatement(decl.Body, fn.Scope); err != nil {
    return nil, err
  }
  return fn.Scope.ReturnValue, nil
}

// func evaluateIndexExpression(exp *ast.IndexExpression, scope *Scope) (ast.Expression, error) {
//  // evaluate the left-hand side
//  left, err := evaluateExpression(exp.Value, scope)
//  if err != nil {
//    return nil, err
//  }
//  indExp, err := evaluateExpression(exp.Index, scope)
//  if err != nil {
//    return nil, err
//  }
//  index, ok := indExp.(*ast.NumberLiteral)
//  if !ok || math.Floor(index.Value) != index.Value || index.Value < 0 {
//    return nil, errors.New("index must evaluate to a whole number >= 0")
//  }
//  // left should be either an array or a string literal
//  i := int(index.Value)
//  if strExp, ok := left.(*ast.StringLiteral); ok {
//    // check array-out-of-bounds
//    if i < 0 || i > len(strExp.Value)-1 {
//      return nil, fmt.Errorf("Index '%d' is out of bounds for string of length %d", i, len(strExp.Value))
//    }
//    return &ast.StringLiteral{
//      Value: string(strExp.Value[i]),
//    }, nil
//  }
//  if arrExp, ok := left.(*ast.ArrayExpression); ok {
//    // check array out of bounds
//    if i < 0 || i > len(arrExp.Expressions)-1 {
//      return nil, fmt.Errorf("Index '%d' is out of bounds for array of length %d", i, len(arrExp.Expressions))
//    }
//    return arrExp.Expressions[i], nil // the array values should already be evaluated
//  }
//  return nil, errors.New("cannot use '@' on non-string or non-array type")
// }
