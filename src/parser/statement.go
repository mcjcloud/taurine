package parser

import (
	"github.com/mcjcloud/taurine/ast"
	"github.com/mcjcloud/taurine/lexer"
	"github.com/mcjcloud/taurine/token"
)

func parseStatement(tkn *token.Token, it *lexer.TokenIterator) ast.Statement {
  if tkn.Type == "{" {
    block := &ast.BlockStatement{Statements: []ast.Statement{}}
    nxt := it.Next()
    for nxt.Type != "}" {
      stmt := parseStatement(nxt, it)
      block.Statements = append(block.Statements, stmt)
      nxt = it.Next()

      if nxt == nil {
        return it.EHandler.Add(nxt, "Expected '}' but found end of file")
      }
    }
    return block
  } else if ast.Symbol(tkn.Value).IsStatementPrefix() {
    if tkn.Value == ast.ETCH {
      return parseEtchStatement(tkn, it)
    } else if tkn.Value == ast.READ {
      return parseReadStatement(tkn, it)
    } else if tkn.Value == ast.IF {
      return parseIfStatement(tkn, it)
    } else if tkn.Value == ast.WHILE {
      return parseWhileLoop(tkn, it)
    } else if tkn.Value == ast.RETURN {
      return parseReturnStatement(tkn, it)
    } else if tkn.Value == ast.IMPORT {
      return parseImportStatement(tkn, it)
    } else if tkn.Value == ast.EXPORT {
      return parseExportStatement(tkn, it)
    }
  } else {
    // it's an expression (symbol)
    exp := parseExpression(tkn, it, nil)
    // expect the semicolon if the expression isn't a block
    if _, ok := exp.(*ast.FunctionLiteral); !ok && it.Next().Type != ";" {
      return it.EHandler.Add(it.Current(), "expected expression statement to end with ';'")
    } else if ok && it.Peek().Type == ";" {
      it.Next()
    }
    return &ast.ExpressionStatement{Expression: exp}
  }
  return it.EHandler.Add(tkn, "unrecognized statement")
}

func parseEtchStatement(tkn *token.Token, it *lexer.TokenIterator) ast.Statement {
  exps := []ast.Expression{}
  nxt := it.Next()
  exp := parseExpression(nxt, it, nil)
  exps = append(exps, exp)
  nxt = it.Next()
  for nxt.Type == "," {
    nxt = it.Next()
    exp = parseExpression(nxt, it, nil)
    exps = append(exps, exp)
    nxt = it.Next()
  }
  if nxt.Type != ";" {
    return it.EHandler.Add(nxt, "expected semicolon to end statement")
  }
  return &ast.EtchStatement{Expressions: exps}
}

func parseReadStatement(tkn *token.Token, it *lexer.TokenIterator) ast.Statement {
  // parse identifier
  nxt := it.Next()
  exp := parseExpression(nxt, it, nil)
  idExp, ok := exp.(*ast.Identifier)
  if !ok {
    return it.EHandler.Add(it.Current(), "expected identifier at beginning of 'read' statement")
  }
  if nxt = it.Next(); nxt.Type != "," && nxt.Type != ";" {
    return it.EHandler.Add(nxt, "expected semicolon to end statement")
  }

  // parse prompt
  exp = parseExpression(it.Next(), it, nil)
  if pmtExp, ok := exp.(*ast.StringLiteral); ok {
    sc := it.Next()
    if sc == nil || sc.Type != ";" {
      return it.EHandler.Add(sc, "expected semicolon to end statement")
    }
    return &ast.ReadStatement{
      Identifier: idExp,
      Prompt:     pmtExp,
    }
  }
  return it.EHandler.Add(it.Current(), "expected prompt after ','")
}

func parseIfStatement(tkn *token.Token, it *lexer.TokenIterator) ast.Statement {
  exp := parseExpression(it.Next(), it, nil)

  stmt := parseStatement(it.Next(), it)
  return &ast.IfStatement{
    Condition: exp,
    Statement: stmt,
  }
}

func parseWhileLoop(tkn *token.Token, it *lexer.TokenIterator) ast.Statement {
  exp := parseExpression(it.Next(), it, nil)

  stmt := parseStatement(it.Next(), it)
  return &ast.WhileLoopStatement{
    Condition: exp,
    Statement: stmt,
  }
}

func parseReturnStatement(tkn *token.Token, it *lexer.TokenIterator) ast.Statement {
  exp := parseExpression(it.Next(), it, nil)
  // expect a semicolon
  if nxt := it.Peek(); nxt.Type != ";" {
    return it.EHandler.Add(it.Current(), "expected semicolon to end return statement")
  }
  it.Next()
  return &ast.ReturnStatement{Value: exp}
}

func parseImportStatement(tkn *token.Token, it *lexer.TokenIterator) ast.Statement {
  ids := make([]*ast.Identifier, 0)
  nxt := it.Next()
  exp := parseExpression(nxt, it, nil)
  if id, ok := exp.(*ast.Identifier); !ok {
    it.SkipStatement()
    return it.EHandler.Add(nxt, "expected identifier.")
  } else {
    ids = append(ids, id)
  }
  for nxt = it.Next(); nxt.Type == ","; nxt = it.Next() {
    idExp := parseExpression(it.Next(), it, nil)
    if id, ok := idExp.(*ast.Identifier); !ok {
      it.SkipStatement()
      return it.EHandler.Add(nxt, "expected identifier.")
    } else {
      ids = append(ids, id)
    }
  }
  // expect FROM
  if nxt.Value != ast.FROM {
    it.SkipStatement()
    return it.EHandler.Add(nxt, "expected 'from'")
  }
  // expect string literal
  if nxt = it.Next(); nxt.Type != "string" {
    it.SkipStatement()
    return it.EHandler.Add(nxt, "expected path to file")
  }
  source := nxt.Value
  // expect semicolon
  if p := it.Peek(); p.Type != ";" {
    return it.EHandler.Add(nxt, "expected ';' to end import statement")
  }
  it.Next()
  return &ast.ImportStatement{
    Source: source,
    Imports: ids,
  }
}

func parseExportStatement(tkn *token.Token, it *lexer.TokenIterator) ast.Statement {
  // parse the exported expression
  exp := parseExpression(it.Next(), it, nil)

  curr := it.Current()
  var nxt *token.Token
  if nxt = it.Next(); nxt.Value == ast.AS {
    // expect an identifier
    idExp := parseExpression(it.Next(), it, nil)
    if id, ok := idExp.(*ast.Identifier); !ok {
      e := it.Current()
      it.SkipStatement()
      return it.EHandler.Add(e, "expected identifier")
    } else {
      return &ast.ExportStatement{
        Identifier: id,
        Value: exp,
      }
    }
  }

  // expect semicolon
  if nxt.Type != ";" {
    return it.EHandler.Add(curr, "expected ';' to end export statement")
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
    return it.EHandler.Add(nxt, "expected variable, function, or identifier")
  }

  // build the export statement
  return &ast.ExportStatement{
    Identifier: id,
    Value: exp,
  }
}

