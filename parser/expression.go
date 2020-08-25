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
		} else if tkn.Type == "symbol" {
			if tkn.Value == "true" {
				return parseExpression(tkn, it, &ast.BooleanLiteral{Value: true})
			} else if tkn.Value == "false" {
				return parseExpression(tkn, it, &ast.BooleanLiteral{Value: false})
			}
			return parseExpression(tkn, it, &ast.Identifier{Name: tkn.Value})
		} else {
			return nil, errors.New("unexpected start of expression")
		}
	}

	// look ahead to see if next token is an operator
	peek := it.Peek()
	if peek.Type == "operation" {
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
	} else if peek.Type == "=" {
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
	}

	if tkn.Type == "number" {
		val, _ := strconv.ParseFloat(tkn.Value, 64)
		return &ast.NumberLiteral{Value: val}, nil
	} else if tkn.Type == "string" {
		return &ast.StringLiteral{Value: tkn.Value}, nil
	} else if tkn.Type == "symbol" {
		if tkn.Value == "true" {
			return &ast.BooleanLiteral{Value: true}, nil
		} else if tkn.Value == "false" {
			return &ast.BooleanLiteral{Value: false}, nil
		}
		return &ast.Identifier{Name: tkn.Value}, nil
	} else {
		return nil, errors.New("unexpected start of expression")
	}
}
