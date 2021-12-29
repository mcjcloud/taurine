package token

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

