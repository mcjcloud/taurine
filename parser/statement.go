package parser

import (
	"errors"

	"github.com/mcjcloud/taurine/ast"
	"github.com/mcjcloud/taurine/lexer"
)

func parseStatement(tkn *lexer.Token, it *lexer.TokenIterator) (ast.Statement, error) {
	if ast.Symbol(tkn.Value).IsStatementPrefix() {
		if tkn.Value == ast.VAR {
			return parseVarDecleration(tkn, it)
		} else if tkn.Value == ast.ETCH {
			return parseEtchStatement(tkn, it)
		} /* else if tkn.Value == ast.FOR {
			return parseForLoop(tkn, it)
		}*/
	}
	return nil, errors.New("Unrecognized statement")
}

func parseVarDecleration(tkn *lexer.Token, it *lexer.TokenIterator) (*ast.VariableDecleration, error) {
	decl := &ast.VariableDecleration{}
	if spec := it.Next(); spec.Type != "(" {
		return nil, errors.New("expected ( after var")
	}

	t := it.Next()
	dataType := ast.Symbol(t.Value)
	if t.Type != "symbol" || !dataType.IsDataType() {
		return nil, errors.New("expected data type after (")
	}
	decl.SymbolType = t.Value

	if spec := it.Next(); spec.Type != ")" {
		return nil, errors.New("expected ) after data type")
	}

	sym := it.Next()
	if sym.Type != "symbol" {
		return nil, errors.New("expected identifier")
	}
	decl.Symbol = sym.Value

	spec := it.Next()
	if spec.Type == "=" {
		// do assignment
		exp := it.Next()
		val, err := parseAssignmentExpression(exp, dataType, it)
		if err != nil {
			return nil, err
		}
		decl.Value = val
		spec = it.Next()
	}
	// TODO: allow multiple assignments with ','
	if spec.Type != ";" {
		return nil, errors.New("missing semicolon")
	}
	return decl, nil
}

func parseAssignmentExpression(tkn *lexer.Token, dataType ast.Symbol, it *lexer.TokenIterator) (ast.Expression, error) {
	exp, err := parseExpression(tkn, it, nil)
	if err != nil {
		return nil, err
	}
	if dataType == ast.NUM {
		if _, ok := exp.(*ast.NumberLiteral); ok {
			return exp, nil
		}
	} else if dataType == ast.STR {
		if _, ok := exp.(*ast.StringLiteral); ok {
			return exp, nil
		}
	} else if dataType == ast.BOOL {
		if _, ok := exp.(*ast.BooleanLiteral); ok {
			return exp, nil
		}
	}
	if _, ok := exp.(*ast.OperationExpression); ok {
		return exp, nil
	}
	return nil, errors.New("assigned type does not match initial value")
}

func parseEtchStatement(tkn *lexer.Token, it *lexer.TokenIterator) (*ast.EtchStatement, error) {
	exps := []ast.Expression{}
	nxt := it.Next()
	exp, err := parseExpression(nxt, it, nil)
	if err != nil {
		return nil, err
	}
	exps = append(exps, exp)
	nxt = it.Next()
	for nxt.Type == "," {
		nxt = it.Next()
		exp, err = parseExpression(nxt, it, nil)
		if err != nil {
			return nil, err
		}
		exps = append(exps, exp)
		nxt = it.Next()
	}
	if nxt.Type != ";" {
		return nil, errors.New("expected semicolon to end statement")
	}
	return &ast.EtchStatement{Expressions: exps}, nil
}

// func parseForLoop(tkn *lexer.Token, it *lexer.TokenIterator) (*ast.ForLoopStatement, error) {

// }
