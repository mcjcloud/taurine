package parser

import (
	"github.com/mcjcloud/taurine/ast"
	"github.com/mcjcloud/taurine/lexer"
)

// Parse parses a series of tokens as a syntax tree
func Parse(tokens []*lexer.Token) (*ast.Source, error) {
	it := lexer.NewTokenIterator(tokens)
	src := &ast.Source{}

	_, tkn := it.Next()
	for tkn != nil {
		if tkn.Type == "symbol" {
			src.BlockStatements = append(src.BlockStatements, parseSymbol(tkn, it))
		}
	}
	return src, nil
}
