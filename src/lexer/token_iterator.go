package lexer

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/mcjcloud/taurine/token"
	"github.com/mcjcloud/taurine/util"
)

// TokenIterator is a structure which allows you to iterate through tokens
type TokenIterator struct {
	Index      int
	Tokens     []*token.Token
  EHandler   *util.ErrorHandler
  IGraph     *util.ImportGraph
  SourcePath string
}

// NewTokenIterator creates a new TokenIterator struct
func NewTokenIterator(tkns []*token.Token, sourcePath string, ig *util.ImportGraph) *TokenIterator {
	return &TokenIterator{
		Index:      -1,
		Tokens:     tkns,
    SourcePath: sourcePath,
    EHandler:   util.NewErrorHandler(),
    IGraph:     ig,
	}
}

func (it *TokenIterator) CreateIteratorForImport(relativePath string) (*TokenIterator, error) {
  absPath := path.Clean(path.Join(path.Dir(it.SourcePath), relativePath))

  // check that the import graph doesn't already contain a parsed AST for this path
  if _, ok := it.IGraph.Nodes[absPath]; ok {
    return nil, &util.AlreadyParsedError{
      Path: absPath,
    }
  }

  // read source code for 
  bytes, err := ioutil.ReadFile(absPath)
  if err != nil {
    return nil, fmt.Errorf("error reading referenced source: %s", err.Error())
  }
  src := string(bytes)

  // tokenize
  tkns := Analyze(src)

  return NewTokenIterator(tkns, absPath, it.IGraph), nil
}

// Peek returns the next token without advancing
func (it *TokenIterator) Peek() *token.Token {
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
func (it *TokenIterator) Next() *token.Token {
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
func (it *TokenIterator) Prev() *token.Token {
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

func PrintTokens(tkns []*token.Token) {
	for i, tkn := range tkns {
    var val string
    if tkn.Type != tkn.Value {
      val = tkn.Value
    }
    fmt.Printf("%02d:%02d %04d %s %s\n", tkn.Position.Row, tkn.Position.Col, i, tkn.Type, val)
	}
}

func (it *TokenIterator) PrintErrors() {
  for _, e := range it.EHandler.Errors {
    // print error message
    fmt.Printf("error at %d:%d: %s\n", e.Token.Position.Row, e.Token.Position.Col, e.Message)

    // print each token in the row with the error
    row := it.GetRow(e.Token.Position.Row)
    colStart := 1
    for _, t := range row {
      // print spaces leading up to the beginning of each token
      for i := colStart; i < t.Position.Col; i += 1 {
        fmt.Printf(" ")
      }
      // update colStart and print the token
      colStart = t.Position.Col+t.Position.Length
      if t.Type == "string" {
        fmt.Printf("\"%s\"", t.Value)
      } else {
        fmt.Printf(t.Value)
      }
    }
    fmt.Println()

    // print underlines up until the error token
    for i := 0; i < e.Token.Position.Col+e.Token.Position.Length-1; i += 1 {
      fmt.Printf("~")
    }
    fmt.Println("^")
  }
}

