package parser

import (
  "errors"
  "strconv"

  "github.com/jinzhu/copier"
  "github.com/mcjcloud/taurine/ast"
  "github.com/mcjcloud/taurine/lexer"
)

func parseExpression(tkn *lexer.Token, it *lexer.TokenIterator, exp ast.Expression) (ast.Expression, error) {
  // if exp is nil, this is the beginning of the expression
  if exp == nil {
    if tkn.Type == "number" {
      val, _ := strconv.ParseFloat(tkn.Value, 64)
      return parseExpression(tkn, it, &ast.NumberLiteral{Value: val})
    } else if tkn.Type == "string" {
      return parseExpression(tkn, it, &ast.StringLiteral{Value: tkn.Value})
    } else if tkn.Type == "bool" {
      // check for boolean value
      if tkn.Value == "true" {
        return parseExpression(tkn, it, &ast.BooleanLiteral{Value: true})
      } else if tkn.Value == "false" {
        return parseExpression(tkn, it, &ast.BooleanLiteral{Value: false})
      }
      return nil, errors.New("invalid boolean value")
    } else if tkn.Type == "symbol" {
     // check if the symbol is "func", if so this is a func expression
      if tkn.Value == "func" {
        fn, err := parseFunction(tkn, it)
        if err != nil {
          return nil, err
        }
        // pass fn back in to see if it is being operated on (e.g. a function call)
        return parseExpression(it.Current(), it, fn)
      }
      return parseExpression(tkn, it, &ast.Identifier{Name: tkn.Value})
    } else if tkn.Type == "[" {
      arrExp, err := parseExpression(it.Next(), it, nil)
      if err != nil {
        return nil, err
      }
      // expect a ]
      nxt := it.Next()
      if nxt == nil || (nxt.Type != "]" && nxt.Type != ",") {
        return nil, errors.New("expected ']' or ',' in array expression")
      }
      exprs := make([]ast.Expression, 1)
      exprs[0] = arrExp
      if nxt.Type == "," {
        // while nxt is a ",", evaluate the next element and add it to the expression array
        for nxt.Type == "," {
          nxtEl, err := parseExpression(it.Next(), it, nil)
          if err != nil {
            return nil, err
          }
          exprs = append(exprs, nxtEl) // add to exp array
          nxt = it.Next()              // get next token
        }
        // check again that it's a closing bracket
        if nxt == nil || nxt.Type != "]" {
          return nil, errors.New("expected ']' to end array expression")
        }
        return parseExpression(nxt, it, &ast.ArrayExpression{Expressions: exprs})
      } else {
        return parseExpression(nxt, it, &ast.ArrayExpression{Expressions: exprs})
      }
    } else if tkn.Type == "(" {
      // (expression)
      grpExp, err := parseExpression(it.Next(), it, nil)
      if err != nil {
        return nil, err
      }
      return parseExpression(it.Next(), it, &ast.GroupExpression{Expression: grpExp})
    } else if tkn.Type == "{" {
      value := make(map[string]ast.Expression)
      keysRemain := true
      var nxt *lexer.Token
      for keysRemain {
        // object literal
        idExp, err := parseExpression(it.Next(), it, nil)
        if err != nil {
          return nil, err
        }
        if id := idExp.(*ast.Identifier); id != nil {
          // expect a ':' next
          if it.Next().Type != ":" {
            return nil, errors.New("expected ':' after identifer")
          }
          valExp, err := parseExpression(it.Next(), it, nil)
          if err != nil {
            return nil, err
          }
          nxt = it.Next()
          if nxt.Type == "," {
            if it.Peek().Type == "}" {
              nxt = it.Next()
              keysRemain = false
            }
          } else if nxt.Type == "}" {
            keysRemain = false
          } else {
            return nil, errors.New("expected ',' or '}' following map key-value pair")
          }
          // add the key value pair to the result
          value[id.Name] = valExp
        } else {
          return nil, errors.New("key must be an identifier")
        }
      }
      return parseExpression(nxt, it, &ast.ObjectLiteral{Value: value})
    } else {
      return nil, errors.New("unexpected start of expression")
    }
  }

  // look ahead to see if next token is an operator
  peek := it.Peek()
  if peek != nil && peek.Type == "operation" {
    op := it.Next()
    rStart := it.Next()
    right, err := parseExpression(rStart, it, nil)
    if err != nil {
      return nil, err
    }
    /*
    if rFnCall, ok := right.(*ast.FunctionCall); ok && op.Type == ast.DOT {
      // TODO: use rFnCall 
    }
    */
    operation := &ast.OperationExpression{
      Operator:        ast.Operator(op.Value),
      LeftExpression:  exp,
      RightExpression: right,
    }
    return orderOperations(operation)
  } else if peek != nil && peek.Type == "=" {
    // assignment
    idExp, ok := exp.(*ast.Identifier)
    if !ok {
      return nil, errors.New("expected left side of assignment to be an identifier")
    }

    it.Next()
    val, err := parseExpression(it.Next(), it, nil)
    if err != nil {
      return nil, err
    }
    return &ast.AssignmentExpression{
      Identifier: idExp,
      Value:      val,
    }, nil
  } else if peek != nil && peek.Type == "(" {
    // function call
    it.Next()
    fnCall, err := parseFunctionCall(exp, it)
    if err != nil {
      return nil, err
    }
    return parseExpression(it.Current(), it, fnCall)
  } else {
    return exp, nil
  }

  /*
  // TODO: consider if tkn.Type needs to be checked. If the above cases don't apply, I think we just return exp
  if tkn.Type == "number" {
    val, _ := strconv.ParseFloat(tkn.Value, 64)
    return &ast.NumberLiteral{Value: val}, nil
  } else if tkn.Type == "string" {
    return &ast.StringLiteral{Value: tkn.Value}, nil
  } else if tkn.Type == "bool" {
    if tkn.Value == "true" {
      return &ast.BooleanLiteral{Value: true}, nil
    } else if tkn.Value == "false" {
      return &ast.BooleanLiteral{Value: false}, nil
    }
    return nil, errors.New("invalid boolean value")
  } else if tkn.Type == "symbol" {
    return &ast.Identifier{Name: tkn.Value}, nil
  } else if grpExp, ok := exp.(*ast.GroupExpression); ok {
    // this ends a group expression
    return grpExp, nil
  } else if arrExp, ok := exp.(*ast.ArrayExpression); ok {
    // this ends an array expression
    return arrExp, nil
  } else if fnExp, ok := exp.(*ast.FunctionCall); ok {
    // this ends a function call expression
    return fnExp, nil
  } else if objExp, ok := exp.(*ast.ObjectLiteral); ok {
    return objExp, nil
  } else if fnExp, ok := exp.(*ast.FunctionLiteral); ok {
    // this ends a function literal
    return fnExp, nil
  } else {
    return nil, errors.New("unexpected start of expression")
  }
  */
}

func orderOperations(opExp *ast.OperationExpression) (*ast.OperationExpression, error) {
  // check if the right child is an operator
  if rightChild, rok := opExp.RightExpression.(*ast.OperationExpression); rok {
    // if so, check the precendence and reorder the tree
    if ast.PRECEDENCE[opExp.Operator] > ast.PRECEDENCE[rightChild.Operator] {
      // copy to avoid modifying the parameter
      opCopy := &ast.OperationExpression{}
      err := copier.Copy(&opCopy, &opExp)
      if err != nil {
        return nil, err
      }

      // set the right child as the new parent and parent as left grandchild
      opCopy.RightExpression = rightChild.LeftExpression
      rightChild.LeftExpression = opCopy

      // recurse to order the right expression
      // TODO: does this need to be done in a loop?
      rightChild.LeftExpression, err = orderOperations(rightChild.LeftExpression.(*ast.OperationExpression))
      if err != nil {
        return nil, err
      }

      // return the new operation
      return rightChild, nil
    }
  }
  return opExp, nil
}

func parseFunction(tkn *lexer.Token, it *lexer.TokenIterator) (*ast.FunctionLiteral, error) {
  // expect ( return type )
  if nxt := it.Next(); nxt == nil || nxt.Type != "(" {
    return nil, errors.New("expected '('")
  }
  nxt := it.Next()
  if nxt == nil || nxt.Type != "symbol" || !ast.Symbol(nxt.Value).IsDataType() {
    return nil, errors.New("expected data type")
  }
  returnType := nxt.Value

  nxt = it.Next()
  if nxt == nil || nxt.Type != ")" {
    return nil, errors.New("expected ')'")
  }

  // expect symbol
  var symbol string
  peek := it.Peek()
  if peek == nil || peek.Type != "symbol" {
    symbol = ""
  } else {
    symbol = it.Next().Value
  }

  // expect ( parameter, parameter, ... )
  params := make([]*ast.VariableDecleration, 0)
  if nxt = it.Next(); nxt == nil || nxt.Type != "(" {
    return nil, errors.New("expected '('")
  }
  for nxt = it.Next(); nxt.Type != ")"; nxt = it.Next() {
    if nxt == nil {
      return nil, errors.New("unexpected end of file")
    }
    if nxt.Type == "," {
      nxt = it.Next()
    }
    // first expect data type
    if !ast.Symbol(nxt.Value).IsDataType() {
      return nil, errors.New("expected data type for parameter")
    }
    dataType := nxt.Value

    // next expect symbol
    nxt = it.Next()
    if nxt == nil || nxt.Type != "symbol" {
      return nil, errors.New("expected parameter name")
    }
    paramName := nxt.Value
    params = append(params, &ast.VariableDecleration{
      Symbol:     paramName,
      SymbolType: dataType,
    })

    // setup for next iteration, should be ',' or ')'
    //peek = it.Peek()
    //if peek == nil || peek.Type != "," && peek.Type != ")" {
    //  return nil, errors.New("expected ')' to end parameters")
    //}
    //if peek.Type == "," {
    //  it.Next()
    //}
  }

  // parse the statement that follows
  body, err := parseStatement(it.Next(), it)
  if err != nil {
    return nil, err
  }
  return &ast.FunctionLiteral{
    Symbol:     symbol,
    ReturnType: returnType,
    Parameters: params,
    Body:       body,
  }, nil
}

func parseFunctionCall(exp ast.Expression, it *lexer.TokenIterator) (*ast.FunctionCall, error) {
  var args []ast.Expression
  nxt := it.Next()
  for nxt.Type != ")" {
    //nxt = it.Next()
    exp, err := parseExpression(nxt, it, nil)
    if err != nil {
      return nil, err
    }
    args = append(args, exp)
    nxt = it.Next()
    if nxt == nil || nxt.Type != "," && nxt.Type != ")" {
      return nil, errors.New("expected ')' to end function call")
    }
  }
  return &ast.FunctionCall{
    Function:  exp,
    Arguments: args,
  }, nil
}
