package lexer

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

// Next advances the iterator by one, returning nil and resetting if the end has been reached
func (it *TokenIterator) Next() *Token {
	it.Index++
	if len(it.Tokens) <= it.Index {
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
