package parser

import (
	"github.com/mcjcloud/taurine/ast"
	"github.com/mcjcloud/taurine/lexer"
)

func parseSymbol(tkn *lexer.Token, it *lexer.TokenIterator) *ast.BlockStatement {
	if tkn.Value == "var" {
		varDecl := parseVarDecleration(tkn, it)
		return &ast.BlockStatement{
			Statements: []*ast.Statement{varDecl},
		}
	}
	panic("oops")
}

func parseVarDecleration(tkn *lexer.Token, it *lexer.TokenIterator) *ast.VariableDecleration {
	decl := &ast.VariableDecleration{}
	_, spec := it.Next()
	if spec.Type != "(" {
		panic(`Expected "(" after var`)
	}
	_, t := it.Next()
	if t.Type != "symbol" {
		panic(`Expected symbol after "("`)
	}
	decl.SymbolType = t.Value
	_, spec = it.Next()
	if spec.Type != ")" {
		panic(`Expected ")" after type`)
	}
	_, sym := it.Next()
	if sym.Type != "symbol" {
		panic(`Expected identifier`)
	}
	decl.Symbol = sym.Value
	_, spec = it.Next()
	if spec.Type == "=" {
		// do assignment
	}
	return decl
}
