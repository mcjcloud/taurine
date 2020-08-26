package parser

import (
	"errors"
	"fmt"

	"github.com/mcjcloud/taurine/ast"
	"github.com/mcjcloud/taurine/lexer"
)

func parseStatement(tkn *lexer.Token, it *lexer.TokenIterator) (ast.Statement, error) {
	if tkn.Type == "{" {
		block := &ast.BlockStatement{Statements: []ast.Statement{}}
		nxt := it.Next()
		for nxt.Type != "}" {
			stmt, err := parseStatement(nxt, it)
			if err != nil {
				return nil, err
			}
			block.Statements = append(block.Statements, stmt)
			nxt = it.Next()

			if nxt == nil {
				return nil, errors.New("Expected '}' but found end of file")
			}
		}
		return block, nil
	} else if ast.Symbol(tkn.Value).IsStatementPrefix() {
		if tkn.Value == ast.VAR {
			return parseVarDecleration(tkn, it)
		} else if tkn.Value == ast.ETCH {
			return parseEtchStatement(tkn, it)
		} else if tkn.Value == ast.READ {
			return parseReadStatement(tkn, it)
		} else if tkn.Value == ast.WHILE {
			return parseWhileLoop(tkn, it)
		}
	} else {
		// it's an expression (identifier)
		exp, err := parseExpression(tkn, it, nil)
		if err != nil {
			return nil, err
		}
		// expect the semicolon
		if it.Next().Type != ";" {
			return nil, errors.New("expected expression statement to end with ';'")
		}
		return &ast.ExpressionStatement{Expression: exp}, nil
	}
	return nil, errors.New("unrecognized statement")
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
	if s := ast.Symbol(sym.Value); s.IsStatementPrefix() || s.IsDataType() {
		return nil, fmt.Errorf("cannot use variable name '%s' as it is a reserved word", s)
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

func parseReadStatement(tkn *lexer.Token, it *lexer.TokenIterator) (*ast.ReadStatement, error) {
	// parse identifier
	nxt := it.Next()
	exp, err := parseExpression(nxt, it, nil)
	if err != nil {
		return nil, err
	}
	idExp, ok := exp.(*ast.Identifier)
	if !ok {
		return nil, errors.New("expected identifier at beginning of 'read' statement")
	}
	if nxt = it.Next(); nxt.Type != "," && nxt.Type != ";" {
		return nil, errors.New("expected semicolon to end statement")
	}

	// parse prompt
	exp, err = parseExpression(it.Next(), it, nil)
	if err != nil {
		return nil, err
	}
	if pmtExp, ok := exp.(*ast.StringLiteral); ok {
		sc := it.Next()
		if sc == nil || sc.Type != ";" {
			return nil, errors.New("expected semicolon to end statement")
		}
		return &ast.ReadStatement{
			Identifier: idExp,
			Prompt:     pmtExp,
		}, nil
	}
	return nil, errors.New("expected prompt after ','")
}

func parseWhileLoop(tkn *lexer.Token, it *lexer.TokenIterator) (*ast.WhileLoopStatement, error) {
	exp, err := parseExpression(it.Next(), it, nil)
	if err != nil {
		return nil, err
	}

	stmt, err := parseStatement(it.Next(), it)
	if err != nil {
		return nil, err
	}
	return &ast.WhileLoopStatement{
		Condition: exp,
		Statement: stmt,
	}, nil
}
