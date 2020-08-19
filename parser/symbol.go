package parser

import (
	"github.com/mcjcloud/taurine/ast"
	"github.com/mcjcloud/taurine/lexer"
)

func parseSymbol(tkn *lexer.Token, it *lexer.TokenIterator) *Node {
	if tkn.Value == "var" {
		return parseVarDecleration(tkn, it)
	}
	panic("oops")
}

func parseVarDecleration(tkn *lexer.Token, it *lexer.TokenIterator) *ast.VariableDecleration {
	if node == nil {
		node = &Node{}
	}
}
