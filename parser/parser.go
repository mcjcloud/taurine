package parser

import (
	"github.com/mcjcloud/taurine/ast"
	"github.com/mcjcloud/taurine/lexer"
)

// Parse parses a series of tokens as a syntax tree
func Parse(tokens []*lexer.Token) (*ast.BlockStatement, error) {
	it := lexer.NewTokenIterator(tokens)
	block := &ast.BlockStatement{}

	tkn := it.Next()
	for tkn != nil {
		if tkn.Type == "{" {
			// block statement

		} else if tkn.Type == "symbol" {
			// statement
			stmt, err := parseStatement(tkn, it)
			if err != nil {
				return nil, err
			}
			block.Statements = append(block.Statements, stmt)
		} else {
			// expression
			exp, err := parseExpression(tkn, it)
			if err != nil {
				return nil, err
			}
			block.Statements = append(block.Statements, &ast.ExpressionStatement{Expression: exp})
		}
		tkn = it.Next()
	}
	return block, nil
}
