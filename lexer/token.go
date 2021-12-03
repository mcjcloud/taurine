package lexer

import "fmt"

// Token represents a token produced by the lexer
type Token struct {
  Type  string
  Value string
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
  return it.Tokens[it.Index+1]
}

// Next advances the iterator by one, returning nil and resetting if the end has been reached
func (it *TokenIterator) Next() *Token {
  it.Index++
  if it.Index >= len(it.Tokens) {
    it.Index = 0
    return nil
  }
  return it.Tokens[it.Index]
}

// Prev moves the iterator back and returns that token
func (it *TokenIterator) Prev() *Token {
  it.Index--
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
    fmt.Printf("%04d %s %s\n", i, tkn.Type, tkn.Value)
  }
}

