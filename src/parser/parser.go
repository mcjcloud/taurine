package parser

import (
  "github.com/mcjcloud/taurine/ast"
  "github.com/mcjcloud/taurine/lexer"
)

// Parse parses a series of tokens as a syntax tree
func Parse(it *lexer.TokenIterator) *ast.Ast {
  block := &ast.BlockStatement{}

  tkn := it.Next()
  for tkn != nil {
    if tkn.Type == "{" || (tkn.Type == "symbol" && ast.Symbol(tkn.Value).IsStatementPrefix()) {
      // statement
      stmt := parseStatement(tkn, it)
      block.Statements = append(block.Statements, stmt)
    } else {
      // expression
      exp := parseExpression(tkn, it, nil)
      // TODO: should probably expect a semicolon here? do some tests.
      block.Statements = append(block.Statements, &ast.ExpressionStatement{Expression: exp})
      // if the expression is not a function, expect an ending semicolon 
      if _, ok := exp.(*ast.FunctionLiteral); !ok {
        errTkn := it.Current()
        if tkn = it.Next(); tkn.Type != ";" {
          it.EHandler.Add(errTkn,  "expected semicolon to end statement")
          continue
        }
      }
    }
    tkn = it.Next()
  }
  return &ast.Ast{
    FilePath:  it.SourcePath,
    Statement: block,
    Exports:   make(map[string]ast.Expression),
  }
}

/*
func JsonAst(stmt *ast.BlockStatement) (string, error)  {
  j, err := json.Marshal(stmt)
  return string(j), err
}
*/

