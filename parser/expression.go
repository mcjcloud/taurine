package parser

import (
	"errors"
	"strconv"

	"github.com/mcjcloud/taurine/ast"
	"github.com/mcjcloud/taurine/lexer"
)

func parseExpression(tkn *lexer.Token, it *lexer.TokenIterator, exp ast.Expression) (ast.Expression, error) {
	// if exp is nil, this is the beginning of the expression
	if exp == nil {
		if tkn.Type == "number" {
			val, _ := strconv.ParseFloat(tkn.Value, 64)
			return parseExpression(tkn, it, &ast.NumberLiteral{Value: val})
		} else if tkn.Type == "string" {
			return parseExpression(tkn, it, &ast.StringLiteral{Value: tkn.Value})
		} else if tkn.Type == "bool" {
			// check for boolean value
			if tkn.Value == "true" {
				return parseExpression(tkn, it, &ast.BooleanLiteral{Value: true})
			} else if tkn.Value == "false" {
				return parseExpression(tkn, it, &ast.BooleanLiteral{Value: false})
			}
			return nil, errors.New("invalid boolean value")
		} else if tkn.Type == "symbol" {
			// check if the identifier is a function call
			if p := it.Peek(); p != nil && p.Type == "(" {
				return parseFunctionCall(tkn, it)
			}
			return parseExpression(tkn, it, &ast.Identifier{Name: tkn.Value})
		} else if tkn.Type == "[" {
			grpExp, err := parseExpression(it.Next(), it, nil)
			if err != nil {
				return nil, err
			}
			// expect a ]
			nxt := it.Next()
			if nxt == nil || nxt.Type != "]" {
				return nil, errors.New("expected ']' to end group expression")
			}
			return parseExpression(nxt, it, &ast.GroupExpression{Expression: grpExp})
		} else {
			return nil, errors.New("unexpected start of expression")
		}
	}

	// look ahead to see if next token is an operator
	peek := it.Peek()
	if peek != nil && peek.Type == "operation" {
		op := it.Next()
		rStart := it.Next()
		right, err := parseExpression(rStart, it, nil)
		if err != nil {
			return nil, err
		}
		return &ast.OperationExpression{
			Operator:        ast.Operator(op.Value),
			LeftExpression:  exp,
			RightExpression: right,
		}, nil
	} else if peek != nil && peek.Type == "=" {
		idExp, ok := exp.(*ast.Identifier)
		if !ok {
			return nil, errors.New("expected left side of assignment to be an identifier")
		}

		it.Next()
		val, err := parseExpression(it.Next(), it, nil)
		if err != nil {
			return nil, err
		}
		return &ast.AssignmentExpression{
			Identifier: idExp,
			Value:      val,
		}, nil
	} else if peek != nil && peek.Type == "@" {
		// the expression here should be a string, array or an identifier of either of those types
		_, sok := exp.(*ast.StringLiteral)
		_, iok := exp.(*ast.Identifier)
		if sok || iok {
			it.Next() // skip the '@'
			index, err := parseExpression(it.Next(), it, nil)
			if err != nil {
				return nil, err
			}
			return &ast.IndexExpression{
				Value: exp,
				Index: index,
			}, nil
		}
		// TODO: else if array
		return nil, errors.New("'@' operator does not apply to non-indexable type")
	}

	if tkn.Type == "number" {
		val, _ := strconv.ParseFloat(tkn.Value, 64)
		return &ast.NumberLiteral{Value: val}, nil
	} else if tkn.Type == "string" {
		return &ast.StringLiteral{Value: tkn.Value}, nil
	} else if tkn.Type == "bool" {
		if tkn.Value == "true" {
			return &ast.BooleanLiteral{Value: true}, nil
		} else if tkn.Value == "false" {
			return &ast.BooleanLiteral{Value: false}, nil
		}
		return nil, errors.New("invalid boolean value")
	} else if tkn.Type == "symbol" {
		return &ast.Identifier{Name: tkn.Value}, nil
	} else if grpExp, ok := exp.(*ast.GroupExpression); ok {
		// this ends a group expression
		return grpExp, nil
	} else {
		return nil, errors.New("unexpected start of expression")
	}
}

func parseFunctionCall(tkn *lexer.Token, it *lexer.TokenIterator) (*ast.FunctionCall, error) {
	var args []ast.Expression
	nxt := it.Next()
	for nxt.Type != ")" {
		nxt = it.Next()
		exp, err := parseExpression(nxt, it, nil)
		if err != nil {
			return nil, err
		}
		args = append(args, exp)
		nxt = it.Next()
		if nxt == nil || nxt.Type != "," && nxt.Type != ")" {
			return nil, errors.New("expected ')' to end function call")
		}
	}
	return &ast.FunctionCall{
		Function:  tkn.Value,
		Arguments: args,
	}, nil
}
