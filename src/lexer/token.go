package lexer

import "fmt"

// Pos represents the position of a Token
type Pos struct {
	Row    int // the row in the file (number of newlines)
	Col    int // the column (number of characters from beginning of row)
	Length int // the length of the token's content
}

// Token represents a token produced by the lexer
type Token struct {
	Type     string
	Value    string
  Position Pos
}

func NewToken(t, v string, scanner Scanner) *Token {
  pos := &Pos{
    Row: scanner.Row,
    Col: scanner.Col-len(v),
    Length: len(v),
  }
  return &Token{
    Type:     t,
    Value:    v,
    Position: *pos,
  }
}

// TokenIterator is a structure which allows you to iterate through tokens
type TokenIterator struct {
	Index  int
	Tokens []*Token
}

// NewTokenIterator creates a new TokenIterator struct
func NewTokenIterator(tkns []*Token) *TokenIterator {
	return &TokenIterator{
		Index:  -1,
		Tokens: tkns,
	}
}

// Peek returns the next token without advancing
func (it *TokenIterator) Peek() *Token {
	if it.Index == len(it.Tokens)-1 {
		return nil
	}
	// find the next non-newline token
	nxt := it.Tokens[it.Index+1]
	for i := 2; nxt.Type == "newline" && it.Index+i < len(it.Tokens); i++ {
		nxt = it.Tokens[it.Index+i]
	}
	return nxt
}

// Next advances the iterator by one, returning nil and resetting if the end has been reached
func (it *TokenIterator) Next() *Token {
	it.Index++
	for it.Index < len(it.Tokens) && it.Tokens[it.Index].Type == "newline" {
		it.Index++
	}
	if it.Index >= len(it.Tokens) {
		it.Index = 0
		return nil
	}
	return it.Tokens[it.Index]
}

// Prev moves the iterator back and returns that token
func (it *TokenIterator) Prev() *Token {
	it.Index--
	for it.Index >= 0 && it.Tokens[it.Index].Type == "newline" {
		it.Index--
	}
	if it.Index < 0 {
		it.Index = 0
		return nil
	}
	return it.Tokens[it.Index]
}

// Current returns the current Token
func (it *TokenIterator) Current() *Token {
	return it.Tokens[it.Index]
}

// AtIndex returns the Token at the given index or nil
func (it *TokenIterator) AtIndex(i int) *Token {
	if i < 0 || i >= len(it.Tokens) {
		return nil
	}
	return it.Tokens[i]
}

func PrintTokens(tkns []*Token) {
	for i, tkn := range tkns {
    var val string
    if tkn.Type != tkn.Value {
      val = tkn.Value
    }
    fmt.Printf("%02d:%02d %04d %s %s\n", tkn.Position.Row, tkn.Position.Col, i, tkn.Type, val)
	}
}
