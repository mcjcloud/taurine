package parser

import (
  "fmt"

	"github.com/mcjcloud/taurine/ast"
	"github.com/mcjcloud/taurine/token"
	"github.com/mcjcloud/taurine/util"
)

func parseStatement(tkn *token.Token, ctx *ParseContext) ast.Statement {
  it := ctx.CurrentIterator()
  if tkn.Type == "{" {
    block := &ast.BlockStatement{Statements: []ast.Statement{}}
    nxt := it.Next()
    for nxt.Type != "}" {
      stmt := parseStatement(nxt, ctx)
      block.Statements = append(block.Statements, stmt)
      nxt = it.Next()

      if nxt == nil {
        return ctx.CurrentErrorHandler().Add(nxt, "Expected '}' but found end of file")
      }
    }
    return block
  } else if ast.Symbol(tkn.Value).IsStatementPrefix() {
    if tkn.Value == ast.ETCH {
      return parseEtchStatement(tkn, ctx)
    } else if tkn.Value == ast.READ {
      return parseReadStatement(tkn, ctx)
    } else if tkn.Value == ast.IF {
      return parseIfStatement(tkn, ctx)
    } else if tkn.Value == ast.FOR {
      return parseForLoop(tkn, ctx)
    } else if tkn.Value == ast.WHILE {
      return parseWhileLoop(tkn, ctx)
    } else if tkn.Value == ast.RETURN {
      return parseReturnStatement(tkn, ctx)
    } else if tkn.Value == ast.IMPORT {
      return parseImportStatement(tkn, ctx)
    } else if tkn.Value == ast.EXPORT {
      return parseExportStatement(tkn, ctx)
    }
  } else {
    // it's an expression (symbol)
    exp := parseExpression(tkn, ctx, nil)
    // expect the semicolon if the expression isn't a block
    if _, ok := exp.(*ast.FunctionLiteral); !ok && it.Next().Type != ";" {
      return ctx.CurrentErrorHandler().Add(it.Current(), "expected expression statement to end with ';'")
    } else if ok && it.Peek().Type == ";" {
      it.Next()
    }
    return &ast.ExpressionStatement{Expression: exp}
  }
  return ctx.CurrentErrorHandler().Add(tkn, "unrecognized statement")
}

func parseEtchStatement(tkn *token.Token, ctx *ParseContext) ast.Statement {
  it := ctx.CurrentIterator()
  exps := []ast.Expression{}
  nxt := it.Next()
  exp := parseExpression(nxt, ctx, nil)
  exps = append(exps, exp)
  nxt = it.Next()
  for nxt.Type == "," {
    nxt = it.Next()
    exp = parseExpression(nxt, ctx, nil)
    exps = append(exps, exp)
    nxt = it.Next()
  }
  if nxt.Type != ";" {
    return ctx.CurrentErrorHandler().Add(nxt, "expected semicolon to end statement")
  }
  return &ast.EtchStatement{Expressions: exps}
}

func parseReadStatement(tkn *token.Token, ctx *ParseContext) ast.Statement {
  it := ctx.CurrentIterator()
  // parse identifier
  nxt := it.Next()
  exp := parseExpression(nxt, ctx, nil)
  idExp, ok := exp.(*ast.Identifier)
  if !ok {
    return ctx.CurrentErrorHandler().Add(it.Current(), "expected identifier at beginning of 'read' statement")
  }
  if nxt = it.Next(); nxt.Type != "," && nxt.Type != ";" {
    return ctx.CurrentErrorHandler().Add(nxt, "expected semicolon to end statement")
  }

  // parse prompt
  exp = parseExpression(it.Next(), ctx, nil)
  if pmtExp, ok := exp.(*ast.StringLiteral); ok {
    sc := it.Next()
    if sc == nil || sc.Type != ";" {
      return ctx.CurrentErrorHandler().Add(sc, "expected semicolon to end statement")
    }
    return &ast.ReadStatement{
      Identifier: idExp,
      Prompt:     pmtExp,
    }
  }
  return ctx.CurrentErrorHandler().Add(it.Current(), "expected prompt after ','")
}

func parseIfStatement(tkn *token.Token, ctx *ParseContext) ast.Statement {
  it := ctx.CurrentIterator()

  exp := parseExpression(it.Next(), ctx, nil)
  stmt := parseStatement(it.Next(), ctx)

  // check for an else [if]
  var elif ast.Statement
  if peek := it.Peek(); peek != nil && peek.Value == ast.ELSE {
    it.Next()
    elif = parseStatement(it.Next(), ctx)
  }

  return &ast.IfStatement{
    Condition: exp,
    Statement: stmt,
    ElseIf:    elif,
  }
}

func parseForLoop(tkn *token.Token, ctx *ParseContext) ast.Statement {
  it := ctx.CurrentIterator()

  // expect an identifier
  idStart := it.Next()
  idExp := parseExpression(idStart, ctx, nil)
  var id *ast.Identifier
  if v, ok := idExp.(*ast.Identifier); !ok {
    it.SkipTo(token.Token{Type: "{", Value: "{"})
    return ctx.CurrentErrorHandler().Add(idStart, fmt.Sprintf("expected identifier but found %s", idExp))
  } else {
    id = v
  }

  // expect 'in' 
  if nxt := it.Next(); nxt.Value != ast.IN {
    it.SkipTo(token.Token{Type: "{", Value: "{"})
    return ctx.CurrentErrorHandler().Add(nxt, fmt.Sprintf("expected 'in' but found %s", nxt.Value))
  }

  // expect expression this should be an array at runtime
  arrExp := parseExpression(it.Next(), ctx, nil)

  // optionally expect a ';' and a number (the step)
  step := 1
  if peek := it.Peek(); peek.Type == ";" {
    s := it.Next()
    numExp := parseExpression(it.Next(), ctx, nil)
    if num, ok := numExp.(*ast.NumberLiteral); ok && num.Value == float64(int(num.Value)) {
      step = int(num.Value)
    } else {
      it.SkipTo(token.Token{Type: "{", Value: "{"})
      return ctx.CurrentErrorHandler().Add(s, fmt.Sprintf("expected integer as step but found %s", numExp))
    }
  }

  // read statement
  stmt := parseStatement(it.Next(), ctx)

  return &ast.ForLoopStatement{
    Control: id,
    Iterator: arrExp,
    Step: step,
    Statement: stmt,
  }
}

func parseWhileLoop(tkn *token.Token, ctx *ParseContext) ast.Statement {
  it := ctx.CurrentIterator()
  exp := parseExpression(it.Next(), ctx, nil)

  stmt := parseStatement(it.Next(), ctx)
  return &ast.WhileLoopStatement{
    Condition: exp,
    Statement: stmt,
  }
}

func parseReturnStatement(tkn *token.Token, ctx *ParseContext) ast.Statement {
  it := ctx.CurrentIterator()
  exp := parseExpression(it.Next(), ctx, nil)
  // expect a semicolon
  if nxt := it.Peek(); nxt.Type != ";" {
    return ctx.CurrentErrorHandler().Add(it.Current(), "expected semicolon to end return statement")
  }
  it.Next()
  return &ast.ReturnStatement{Value: exp}
}

func parseImportStatement(tkn *token.Token, ctx *ParseContext) ast.Statement {
  it := ctx.CurrentIterator()
  handler := ctx.CurrentErrorHandler()
  ids := make([]*ast.Identifier, 0)
  nxt := it.Next()
  exp := parseExpression(nxt, ctx, nil)
  if id, ok := exp.(*ast.Identifier); !ok {
    it.SkipStatement()
    return handler.Add(nxt, "expected identifier.")
  } else {
    ids = append(ids, id)
  }
  for nxt = it.Next(); nxt.Type == ","; nxt = it.Next() {
    idExp := parseExpression(it.Next(), ctx, nil)
    if id, ok := idExp.(*ast.Identifier); !ok {
      it.SkipStatement()
      return handler.Add(nxt, "expected identifier.")
    } else {
      ids = append(ids, id)
    }
  }
  // expect FROM
  if nxt.Value != ast.FROM {
    it.SkipStatement()
    return handler.Add(nxt, "expected 'from'")
  }
  // expect string literal
  if nxt = it.Next(); nxt.Type != "string" {
    it.SkipStatement()
    return handler.Add(nxt, "expected path to file")
  }
  source := nxt.Value
  // expect semicolon
  if p := it.Peek(); p.Type != ";" {
    return handler.Add(nxt, "expected ';' to end import statement")
  }
  it.Next()

  // PushImport updates the context to start parsing the referenced file
  err := ctx.PushImport(source)
  if _, ok := err.(*util.AlreadyParsedError); !ok && err != nil {
    return handler.Add(nxt, "error finding referenced file")
  } else if ok {
    return &ast.ImportStatement{
      Source: source,
      Imports: ids,
    }
  }

  // run Parse and then return ctx to previous state
  refTree := Parse(ctx)
  ctx.PopImportWithTree(refTree)

  // return the import statement node
  return &ast.ImportStatement{
    Source: source,
    Imports: ids,
  }
}

func parseExportStatement(tkn *token.Token, ctx *ParseContext) ast.Statement {
  it := ctx.CurrentIterator()
  // parse the exported expression
  exp := parseExpression(it.Next(), ctx, nil)

  curr := it.Current()
  var nxt *token.Token
  if nxt = it.Next(); nxt.Value == ast.AS {
    // expect an identifier
    idExp := parseExpression(it.Next(), ctx, nil)
    if id, ok := idExp.(*ast.Identifier); !ok {
      e := it.Current()
      it.SkipStatement()
      return ctx.CurrentErrorHandler().Add(e, "expected identifier")
    } else {
      return &ast.ExportStatement{
        Identifier: id,
        Value: exp,
      }
    }
  }

  // expect semicolon
  if _, ok := exp.(*ast.FunctionLiteral); !ok && nxt.Type != ";" {
    return ctx.CurrentErrorHandler().Add(curr, "expected ';' to end export statement")
  } else if nxt.Type != ";" {
    it.Prev()
  }

  // if AS is not used, the identifier should exist in the value
  // either a function or variable
  var id *ast.Identifier
  if fn, ok := exp.(*ast.FunctionLiteral); ok {
    id = &ast.Identifier{
      Name: fn.Symbol,
    }
  } else if v, ok := exp.(*ast.VariableDecleration); ok {
    id = &ast.Identifier{
      Name: v.Symbol,
    }
  } else if i, ok := exp.(*ast.Identifier); ok {
    id = &ast.Identifier{
      Name: i.Name,
    }
  } else {
    return ctx.CurrentErrorHandler().Add(nxt, "expected variable, function, or identifier")
  }

  // build the export statement
  return &ast.ExportStatement{
    Identifier: id,
    Value: exp,
  }
}

