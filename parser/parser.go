package parser

import (
  "encoding/json"

  "github.com/mcjcloud/taurine/ast"
  "github.com/mcjcloud/taurine/lexer"
)

// Parse parses a series of tokens as a syntax tree
func Parse(tokens []*lexer.Token) (*ast.BlockStatement, error) {
  it := lexer.NewTokenIterator(tokens)
  block := &ast.BlockStatement{}

  tkn := it.Next()
  for tkn != nil {
    if tkn.Type == "{" || tkn.Type == "symbol" {
      // statement
      stmt, err := parseStatement(tkn, it)
      if err != nil {
        return nil, err
      }
      block.Statements = append(block.Statements, stmt)
    } else {
      // expression
      exp, err := parseExpression(tkn, it, nil)
      if err != nil {
        return nil, err
      }
      block.Statements = append(block.Statements, &ast.ExpressionStatement{Expression: exp})
    }
    tkn = it.Next()
  }
  return block, nil
}

func JsonAst(stmt *ast.BlockStatement) (string, error)  {
  j, err := json.Marshal(stmt)
  return string(j), err
}

