package lexer

import (
	"fmt"

	"github.com/mcjcloud/taurine/token"
)

// TokenIterator is a structure which allows you to iterate through tokens
type TokenIterator struct {
	Index      int
	Tokens     []*token.Token
}

// NewTokenIterator creates a new TokenIterator struct
func NewTokenIterator(tkns []*token.Token) *TokenIterator {
	return &TokenIterator{
		Index:      -1,
		Tokens:     tkns,
	}
}

// Peek returns the next token without advancing
func (it *TokenIterator) Peek() *token.Token {
	if it.Index == len(it.Tokens)-1 {
		return nil
	}
	// find the next non-newline token
  var i int
	for i = 1; it.Index+i < len(it.Tokens) && it.Tokens[it.Index+i].Type == "newline"; i++ {
    continue
	}
  if it.Index+i < len(it.Tokens) {
    return it.Tokens[it.Index+i]
  }
	return nil
}

// Next advances the iterator by one, returning nil and resetting if the end has been reached
func (it *TokenIterator) Next() *token.Token {
	it.Index++
	for it.Index < len(it.Tokens) && it.Tokens[it.Index].Type == "newline" {
		it.Index++
	}
	if it.Index >= len(it.Tokens) {
		return nil
	}
	return it.Tokens[it.Index]
}

// Prev moves the iterator back and returns that token
func (it *TokenIterator) Prev() *token.Token {
	it.Index--
	for it.Index >= 0 && it.Tokens[it.Index].Type == "newline" {
		it.Index--
	}
	if it.Index < 0 {
		return nil
	}
	return it.Tokens[it.Index]
}

// Current returns the current Token
func (it *TokenIterator) Current() *token.Token {
	return it.Tokens[it.Index]
}

// AtIndex returns the Token at the given index or nil
func (it *TokenIterator) AtIndex(i int) *token.Token {
	if i < 0 || i >= len(it.Tokens) {
		return nil
	}
	return it.Tokens[i]
}

// SkipStatement advance the iterator past the next ';' or '}'
func (it *TokenIterator) SkipStatement() {
  for nxt := it.Next(); nxt.Type != ";"; nxt = it.Next() {}
}

// SkipToClosingBracket skips to the closing bracket matching the last '{'
func (it *TokenIterator) SkipToClosingBracket() {
  var depth int
  for nxt := it.Next(); nxt.Type != "}" || depth > 0; nxt = it.Next() {
    if nxt.Type == "{" {
      depth += 1
    } else if nxt.Type == "}" {
      depth -= 1
    }
  }
}

// SkipTo advances the iterator to the next occurance of the given token
func (it *TokenIterator) SkipTo(tkn token.Token) {
  for nxt := it.Next(); nxt.Type != tkn.Type || nxt.Value != tkn.Value; nxt = it.Next() {}
}

// GetRow returns all the tokens on the given row
func (it *TokenIterator) GetRow(n int) []*token.Token {
  row := make([]*token.Token, 0)
  for _, t := range it.Tokens {
    if t.Position.Row < n {
      continue
    }
    if t.Position.Row > n {
      break
    }
    row = append(row, t)
  }
  return row
}

func (it TokenIterator) PrintTokens() {
	for i, tkn := range it.Tokens {
    var val string
    if tkn.Type != tkn.Value {
      val = tkn.Value
    }
    fmt.Printf("%02d:%02d %04d %s %s\n", tkn.Position.Row, tkn.Position.Col, i, tkn.Type, val)
	}
}

