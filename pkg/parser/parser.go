package parser

import (
	"github.com/mcjcloud/taurine/pkg/ast"
)

// Parse parses a series of tokens as a syntax tree
func Parse(ctx *ParseContext) *ast.Ast {
	it := ctx.CurrentIterator()
	handler := ctx.CurrentErrorHandler()
	block := &ast.BlockStatement{}

	tkn := it.Next()
	for tkn != nil {
		if tkn.Type == "{" || (tkn.Type == "symbol" && ast.Symbol(tkn.Value).IsStatementPrefix()) {
			// statement
			stmt := parseStatement(tkn, ctx)
			block.Statements = append(block.Statements, stmt)
		} else {
			// expression
			exp := parseExpression(tkn, ctx, nil)
			// TODO: should probably expect a semicolon here? do some tests.
			block.Statements = append(block.Statements, &ast.ExpressionStatement{Expression: exp})
			// if the expression is not a function, expect an ending semicolon
			if _, ok := exp.(*ast.FunctionLiteral); !ok {
				errTkn := it.Current()
				if tkn = it.Next(); tkn.Type != ";" {
					handler.Add(errTkn, "expected semicolon to end statement")
					continue
				}
			}
		}
		tkn = it.Next()
	}
	return &ast.Ast{
		FilePath:  ctx.CurrentFilePath(),
		Statement: block,
		Exports:   make(map[string]ast.Expression),
	}
}
