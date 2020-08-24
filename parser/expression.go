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
			numExp := &ast.NumberLiteral{Value: val}
			return parseExpression(tkn, it, numExp)
		} else if tkn.Type == "string" {
			strExp := &ast.StringLiteral{Value: tkn.Value}
			return parseExpression(tkn, it, strExp)
		} else if tkn.Type == "symbol" {
			idExp := &ast.Identifier{Name: tkn.Value}
			return parseExpression(tkn, it, idExp)
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
	}
	if tkn.Type == "number" {
		val, _ := strconv.ParseFloat(tkn.Value, 64)
		return &ast.NumberLiteral{Value: val}, nil
	} else if tkn.Type == "string" {
		return &ast.StringLiteral{Value: tkn.Value}, nil
	} else if tkn.Type == "symbol" {
		return &ast.Identifier{Name: tkn.Value}, nil
	} else {
		return nil, errors.New("unexpected start of expression")
	}
}

// func parseOperation(left *ast.ExpressionStatement, op ast.Operator, it *lexer.TokenIterator) (*ast.OperationExpression, error) {
// 	nxt := it.Next()
// 	exp, err := parseExpression(nxt, it, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &ast.OperationExpression{
// 		LeftExpression:  left.Expression,
// 		Operator:        op,
// 		RightExpression: exp,
// 	}, nil
// }
