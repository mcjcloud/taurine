package parser

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/mcjcloud/taurine/ast"
	"github.com/mcjcloud/taurine/token"
)

func parseExpression(tkn *token.Token, ctx *ParseContext, exp ast.Expression) ast.Expression {
  it := ctx.CurrentIterator()
  if exp == nil {
    if tkn.Type == "number" {
      val, _ := strconv.ParseFloat(tkn.Value, 64)
      return parseExpression(tkn, ctx, &ast.NumberLiteral{Value: val})
    } else if tkn.Type == "integer" {
      bigInt, _ := new(big.Int).SetString(tkn.Value, 10)
      return parseExpression(tkn, ctx, &ast.IntegerLiteral{Value: bigInt})
    } else if tkn.Type == "string" {
      return parseExpression(tkn, ctx, &ast.StringLiteral{Value: tkn.Value})
    } else if tkn.Type == "bool" {
      // check for boolean value
      if tkn.Value == "true" {
        return parseExpression(tkn, ctx, &ast.BooleanLiteral{Value: true})
      } else if tkn.Value == "false" {
        return parseExpression(tkn, ctx, &ast.BooleanLiteral{Value: false})
      }
      return ctx.CurrentErrorHandler().Add(tkn, "invalid boolean value")
    } else if tkn.Type == "symbol" {
      // check if the symbol is "func", if so this is a func expression
      if tkn.Value == ast.FUNC {
        fn := parseFunction(tkn, ctx)
        // pass fn back in to see if it is being operated on (e.g. a function call)
        return parseExpression(it.Current(), ctx, fn)
      } else if tkn.Value == ast.VAR {
        vDecl := parseVarDeclaration(tkn, ctx)
        return parseExpression(it.Current(), ctx, vDecl)
      } else {
        return parseExpression(tkn, ctx, &ast.Identifier{Name: tkn.Value})
      }
    } else if tkn.Type == "[" {
      arrExp := parseExpression(it.Next(), ctx, nil)
      // expect a ]
      nxt := it.Next()
      if nxt == nil || (nxt.Type != "]" && nxt.Type != ",") {
        it.SkipStatement()
        return ctx.CurrentErrorHandler().Add(nxt, "expected ']' or ',' in array expression")
      }
      exprs := make([]ast.Expression, 1)
      exprs[0] = arrExp
      if nxt.Type == "," {
        // while nxt is a ",", evaluate the next element and add it to the expression array
        for nxt.Type == "," {
          nxtEl := parseExpression(it.Next(), ctx, nil)
          exprs = append(exprs, nxtEl) // add to exp array
          nxt = it.Next()              // get next token
        }
        // check again that it's a closing bracket
        if nxt == nil || nxt.Type != "]" {
          it.SkipStatement()
          return ctx.CurrentErrorHandler().Add(nxt, "expected ']' to end array expression")
        }
        return parseExpression(nxt, ctx, &ast.ArrayExpression{Expressions: exprs})
      } else {
        return parseExpression(nxt, ctx, &ast.ArrayExpression{Expressions: exprs})
      }
    } else if tkn.Type == "(" {
      // (expression)
      grpExp := parseExpression(it.Next(), ctx, nil)
      return parseExpression(it.Next(), ctx, &ast.GroupExpression{Expression: grpExp})
    } else if tkn.Type == "{" {
      // object
      value := make(map[string]ast.Expression)
      keysRemain := true
      var nxt *token.Token
      for keysRemain {
        // object literal
        idExp := parseExpression(it.Next(), ctx, nil)
        if id := idExp.(*ast.Identifier); id != nil {
          // expect a ':' next
          if it.Next().Type != ":" {
            it.SkipToClosingBracket()
            return ctx.CurrentErrorHandler().Add(it.Current(), "expected ':' after identifer")
          }
          valExp := parseExpression(it.Next(), ctx, nil)
          nxt = it.Next()
          if nxt.Type == "," {
            if it.Peek().Type == "}" {
              nxt = it.Next()
              keysRemain = false
            }
          } else if nxt.Type == "}" {
            keysRemain = false
          } else {
            it.SkipToClosingBracket()
            return ctx.CurrentErrorHandler().Add(nxt, "expected ',' or '}' following map key-value pair")
          }
          // add the key value pair to the result
          value[id.Name] = valExp
        } else {
          // skip to the closing bracket
          it.SkipToClosingBracket()
          return ctx.CurrentErrorHandler().Add(it.Current(), "key must be an identifier")
        }
      }
      return parseExpression(nxt, ctx, &ast.ObjectLiteral{Value: value})
    } else {
      return ctx.CurrentErrorHandler().Add(tkn, fmt.Sprintf("unexpected start of expression: (%d, %s)", ctx.CurrentIterator().Index, tkn.Type))
    }
  }

  // look ahead to see if next token is an operator
  peek := it.Peek()
  if peek != nil && peek.Type == "operation" {
    op := it.Next()
    rStart := it.Next()
    right := parseExpression(rStart, ctx, nil)
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
      return ctx.CurrentErrorHandler().Add(peek, "expected left side of assignment to be an identifier")
    }

    it.Next()
    val := parseExpression(it.Next(), ctx, nil)
    return &ast.AssignmentExpression{
      Identifier: idExp,
      Value:      val,
    }
  } else if peek != nil && peek.Type == "(" {
    // function call
    it.Next()
    fnCall := parseFunctionCall(exp, ctx)
    return parseExpression(it.Current(), ctx, fnCall)
  } else {
    return exp
  }
}

func orderOperations(opExp *ast.OperationExpression) *ast.OperationExpression {
  // check if the right child is an operator
  if rightChild, rok := opExp.RightExpression.(*ast.OperationExpression); rok {
    // if so, check the precendence and reorder the tree
    if ast.PRECEDENCE[opExp.Operator] > ast.PRECEDENCE[rightChild.Operator] {
      // copy to avoid modifying the parameter
      opCopy := &ast.OperationExpression{}
      err := copier.Copy(&opCopy, &opExp)
      if err != nil {
        panic(err)
      }

      // set the right child as the new parent and parent as left grandchild
      opCopy.RightExpression = rightChild.LeftExpression
      rightChild.LeftExpression = opCopy

      // recurse to order the right expression
      rightChild.LeftExpression = orderOperations(rightChild.LeftExpression.(*ast.OperationExpression))

      // return the new operation
      return rightChild
    }
  }
  return opExp
}

func parseVarDeclaration(tkn *token.Token, ctx *ParseContext) ast.Expression {
  it := ctx.CurrentIterator()
  decl := &ast.VariableDecleration{}
  if spec := it.Next(); spec.Type != "(" {
    it.SkipStatement()
    return ctx.CurrentErrorHandler().Add(spec, "expected '(' after var")
  }

  t := it.Next()
  dataType := ast.Symbol(t.Value)
  if t.Type != "symbol" || !dataType.IsDataType() {
    it.SkipStatement()
    return ctx.CurrentErrorHandler().Add(t, "expected data type after (")
  }
  decl.SymbolType = t.Value

  if spec := it.Next(); spec.Type != ")" {
    it.SkipStatement()
    return ctx.CurrentErrorHandler().Add(spec, "expected ) after data type")
  }

  sym := it.Next()
  if sym.Type != "symbol" {
    it.SkipStatement()
    return ctx.CurrentErrorHandler().Add(sym, "expected identifier")
  }
  // TODO: this won't work properly. Create another method for reserved words
  if s := ast.Symbol(sym.Value); s.IsStatementPrefix() || s.IsDataType() {
    it.SkipStatement()
    return ctx.CurrentErrorHandler().Add(sym, fmt.Sprintf("cannot use variable name '%s' as it is a reserved word", s))
  }
  decl.Symbol = sym.Value

  spec := it.Next()
  if spec.Type == "=" {
    // do assignment
    exp := it.Next()
    val := parseAssignmentExpression(exp, dataType, ctx)
    decl.Value = val
  } else {
    it.Prev()
  }
  // TODO: allow multiple assignments with ','
  return decl
}

func parseFunction(tkn *token.Token, ctx *ParseContext) ast.Expression {
  it := ctx.CurrentIterator()
  // expect ( return type )
  if nxt := it.Next(); nxt == nil || nxt.Type != "(" {
    // skip to the opening bracket, and then the closing one
    it.SkipTo(token.Token{Type: "{", Value: "{"})
    it.SkipToClosingBracket()
    return ctx.CurrentErrorHandler().Add(tkn, "expected '('")
  }
  nxt := it.Next()
  if nxt == nil || nxt.Type != "symbol" || !ast.Symbol(nxt.Value).IsDataType() {
    // skip to the opening bracket, and then the closing one
    it.SkipTo(token.Token{Type: "{", Value: "{"})
    it.SkipToClosingBracket()
    return ctx.CurrentErrorHandler().Add(nxt, "expected data type")
  }
  returnType := nxt.Value

  nxt = it.Next()
  if nxt == nil || nxt.Type != ")" {
    // skip to the opening bracket, and then the closing one
    it.SkipTo(token.Token{Type: "{", Value: "{"})
    it.SkipToClosingBracket()
    return ctx.CurrentErrorHandler().Add(nxt, "expected ')'")
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
    // skip to the opening bracket, and then the closing one
    it.SkipTo(token.Token{Type: "{", Value: "{"})
    it.SkipToClosingBracket()
    return ctx.CurrentErrorHandler().Add(tkn, "expected '('")
  }
  for nxt = it.Next(); nxt.Type != ")"; nxt = it.Next() {
    if nxt == nil {
      return ctx.CurrentErrorHandler().Add(tkn, "unexpected end of file")
    }
    if nxt.Type == "," {
      nxt = it.Next()
    }
    // first expect data type
    if !ast.Symbol(nxt.Value).IsDataType() {
      // skip to the opening bracket, and then the closing one
      it.SkipTo(token.Token{Type: "{", Value: "{"})
      it.SkipToClosingBracket()
      return ctx.CurrentErrorHandler().Add(nxt, "expected data type for parameter")
    }
    dataType := nxt.Value

    // next expect symbol
    nxt = it.Next()
    if nxt == nil || nxt.Type != "symbol" {
      // skip to the opening bracket, and then the closing one
      it.SkipTo(token.Token{Type: "{", Value: "{"})
      it.SkipToClosingBracket()
      return ctx.CurrentErrorHandler().Add(nxt, "expected parameter name")
    }
    paramName := nxt.Value
    params = append(params, &ast.VariableDecleration{
      Symbol:     paramName,
      SymbolType: dataType,
    })
  }

  // parse the statement that follows
  body := parseStatement(it.Next(), ctx)
  return &ast.FunctionLiteral{
    Symbol:     symbol,
    ReturnType: returnType,
    Parameters: params,
    Body:       body,
  }
}

func parseFunctionCall(exp ast.Expression, ctx *ParseContext) ast.Expression {
  it := ctx.CurrentIterator()
  var args []ast.Expression
  nxt := it.Next()
  for nxt.Type != ")" {
    exp := parseExpression(nxt, ctx, nil)
    args = append(args, exp)
    nxt = it.Next()
    if nxt == nil || nxt.Type != "," && nxt.Type != ")" {
      it.SkipTo(token.Token{Type: ")", Value: ")"})
      return ctx.CurrentErrorHandler().Add(nxt, "expected ')' to end function call")
    }
    if nxt.Type == "," {
      nxt = it.Next()
    }
  }
  return &ast.FunctionCall{
    Function:  exp,
    Arguments: args,
  }
}

func parseAssignmentExpression(tkn *token.Token, dataType ast.Symbol, ctx *ParseContext) ast.Expression {
  exp := parseExpression(tkn, ctx, nil)
  if dataType == ast.NUM {
    if _, ok := exp.(*ast.NumberLiteral); ok {
      return exp
    } else if intLit, ok := exp.(*ast.IntegerLiteral); ok {
      return &ast.NumberLiteral{Value: float64(intLit.Value.Int64())}
    }
  } else if dataType == ast.INT {
    if _, ok := exp.(*ast.IntegerLiteral); ok {
      return exp
    }
  } else if dataType == ast.STR {
    if _, ok := exp.(*ast.StringLiteral); ok {
      return exp
    }
  } else if dataType == ast.BOOL {
    if _, ok := exp.(*ast.BooleanLiteral); ok {
      return exp
    }
  } else if dataType == ast.ARR {
    if _, ok := exp.(*ast.ArrayExpression); ok {
      return exp
    }
  } else if dataType == ast.OBJ {
    if _, ok := exp.(*ast.ObjectLiteral); ok {
      return exp
    }
  } else if dataType == ast.FUNC {
    if _, ok := exp.(*ast.FunctionLiteral); ok {
      return exp
    }
  }
  if _, ok := exp.(*ast.OperationExpression); ok {
    return exp
  }
  if _, ok := exp.(*ast.FunctionCall); ok {
    return exp
  }
  if _, ok := exp.(*ast.GroupExpression); ok {
    return exp
  }
  if _, ok := exp.(*ast.FunctionCall); ok {
    return exp
  }
  return ctx.CurrentErrorHandler().Add(ctx.CurrentIterator().Current(), "assigned type does not match initial value")
}

