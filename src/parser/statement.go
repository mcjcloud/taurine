package parser

import (
  "github.com/mcjcloud/taurine/ast"
  "github.com/mcjcloud/taurine/token"
  "github.com/mcjcloud/taurine/lexer"
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
