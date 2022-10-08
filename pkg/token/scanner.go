package token

import (
  "io"
	"strings"
)

const EOF byte = 0

// Scanner is used by a lexer to keep track of source code positions
type Scanner struct {
  // inherits from Reader
  *strings.Reader
  rowLengths   []int

  Source       string
  SourceLength int
  Row          int
  Col          int
}

func NewScanner(src string) *Scanner {
  return &Scanner{
    Reader: strings.NewReader(src),
    Source:       src,
    SourceLength: len(src),
    Row:          1,
    Col:          1,
    rowLengths:   []int{1},
  }
}

func (s *Scanner) HasNext() bool {
  return s.Len() > 0
}

func (s *Scanner) Next() byte {
  c, err := s.ReadByte()
  if err != nil {
    if err == io.EOF {
      return EOF
    }
    panic(err)
  }
  if c == '\n' {
    s.Row += 1
    s.Col = 1
    if len(s.rowLengths) < s.Row {
      s.rowLengths = append(s.rowLengths, 1)
    }
  } else {
    s.Col += 1
    s.rowLengths[s.Row -1] = s.Col
  }
  return c
}

func (s *Scanner) Unread() {
  err := s.UnreadByte()
  if err != nil {
    panic(err)
  }
  s.Col -= 1
  if s.Col <= 0 {
    s.Row -= 1
    s.Col = s.rowLengths[s.Row-1]
  }
}

