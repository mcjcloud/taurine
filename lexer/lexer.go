package lexer

import (
  "fmt"
  "io"
  "regexp"
  "strings"
)

var numberRe = regexp.MustCompile(`[.0-9]`)
var symbolRe = regexp.MustCompile(`[_a-zA-Z0-9]`)
var boolRe = regexp.MustCompile(`(^true$)|(^false$)`)

func isWhitespace(c byte) bool {
  return c == ' ' || c == '\n' || c == '\t' || c == '\r'
}

func isSpecial(c byte) bool {
  return c == '{' || c == '}' || c == '(' || c == ')' || c == ',' || c == ';' || c == ':' || c == '[' || c == ']'
}

func isOperation(c byte) bool {
  return c == '+' || c == '-' || c == '*' || c == '/' || c == '%' || c == '!' || c == '<' || c == '>' || c == '@' || c == '.'
}

// Analyze creates a series of tokens from source code
func Analyze(source string) (tkns []*Token) {
  tkns = make([]*Token, 0)
  srcReader := strings.NewReader(source)
  for srcReader.Len() > 0 {
    c, err := srcReader.ReadByte()
    if err != nil {
      panic(err)
    }

    // skip whitespace
    if isWhitespace(c) {
      continue
    }

    if c == '/' {
      // check if the next character is a /
      nxt, err := srcReader.ReadByte()
      if err != nil {
        panic(err)
      }
      if nxt == '/' {
        // eat every character until a newline is found
        for nxt != '\n' && nxt != '\r' {
          nxt, err = srcReader.ReadByte()
          if err != nil {
            panic(err)
          }
        }
        continue
      } else {
        // otherwise put both of the characters back and keep going
        srcReader.UnreadByte()
      }
    }

    if c == '"' {
      tkns = append(tkns, scanString(srcReader))
    } else if c == '-' {
      nxt, err := srcReader.ReadByte()
      if err != nil {
        panic(err)
      }
      if numberRe.Match([]byte{nxt}) {
        err = srcReader.UnreadByte()
        if err != nil {
          panic(err)
        }
        tkns = append(tkns, scanNumber(c, srcReader))
      }
    } else if isSpecial(c) {
      tkns = append(tkns, &Token{Type: string(c)}) // special characters {}()@,;:= will be their own type
    } else if isOperation(c) {
      tkns = append(tkns, scanOperation(c, srcReader))
    } else if c == '=' {
      // check the next one is '=' to see if this is special or an operation
      nxt, err := srcReader.ReadByte()
      if err != nil {
        panic(err)
      }
      if nxt == '=' {
        tkns = append(tkns, &Token{Type: "operation", Value: "=="})
      } else {
        if err = srcReader.UnreadByte(); err != nil {
          panic(err)
        }
        tkns = append(tkns, &Token{Type: string(c)})
      }
    } else if numberRe.Match([]byte{c}) { // number literal
      tkns = append(tkns, scanNumber(c, srcReader))
    } else if symbolRe.Match(([]byte{c})) { // symbol
      tkn := scan(c, srcReader, symbolRe, "symbol")
      if boolRe.MatchString(tkn.Value) { // boolean
        tkns = append(tkns, &Token{Type: "bool", Value: tkn.Value})
      } else {
        tkns = append(tkns, tkn)
      }
    } else {
      panic(fmt.Sprintf("Unexpected character %v\n", c))
    }
  }
  return
}

// scan a string from the reader, including the double quotes
func scanString(reader *strings.Reader) *Token {
  var val string
  c, err := reader.ReadByte()
  if err != nil {
    panic(err)
  }
  for c != '"' {
    val += string(c)
    c, err = reader.ReadByte()
    if err != nil {
      panic(err)
    }
    if c == '\\' {
      val += string(c)
      c, err = reader.ReadByte()
      if err != nil {
        panic(err)
      }
      val += string(c)
      c, err = reader.ReadByte()
      if err != nil {
        panic(err)
      }
    }
  }
  return &Token{
    Type:  "string",
    Value: val,
  }
}

func scanOperation(c byte, reader *strings.Reader) *Token {
  val := string(c)
  b, err := reader.ReadByte()
  if err != nil {
    if err == io.EOF {
      return &Token{
        Type:  "operation",
        Value: val,
      }
    }
    panic(err)
  }

  if b == '=' {
    val += string(b)
  } else {
    if err = reader.UnreadByte(); err != nil {
      panic(err)
    }
  }
  return &Token{
    Type:  "operation",
    Value: val,
  }
}

func scanNumber(c byte, reader *strings.Reader) *Token {
  var val string
  var err error
  b := c
  if c == '-' {
    val += "-"
    b, err = reader.ReadByte()
    if err != nil {
      panic(err)
    }
  }
  for numberRe.Match([]byte{b}) {
    val += string(b)
    b, err = reader.ReadByte()
    if err != nil {
      if err == io.EOF {
        break
      }
      panic(err)
    }
  }
  if err == nil {
    err = reader.UnreadByte()
  }
  if err != nil && err != io.EOF {
    panic(err)
  }
  return &Token{
    Type:  "number",
    Value: val,
  }
}

func scan(c byte, reader *strings.Reader, re *regexp.Regexp, t string) *Token {
  var val string
  var err error
  b := c
  for re.Match([]byte{b}) {
    val += string(b)
    b, err = reader.ReadByte()
    if err != nil {
      if err == io.EOF {
        break
      }
      panic(err)
    }
  }
  if err == nil {
    err = reader.UnreadByte()
  }
  if err != nil && err != io.EOF {
    panic(err)
  }
  return &Token{
    Type:  t,
    Value: val,
  }
}

