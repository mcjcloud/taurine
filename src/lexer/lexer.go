package lexer

import (
  "fmt"
  "regexp"

  "github.com/mcjcloud/taurine/token"
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
func Analyze(source string) (tkns []*token.Token) {
  tkns = make([]*token.Token, 0)
  scanner := token.NewScanner(source)

  // vars for tracking token positions
  for scanner.HasNext() {
    c := scanner.Next()

    // skip whitespace
    if isWhitespace(c) {
      if c == '\n' {
        //tkns = append(tkns, &Token{Type: "newline"})
        tkns = append(tkns, token.NewToken("newline", "", *scanner))
      }
      continue

    }
    if c == '/' {
      // check if the next character is a /
      nxt := scanner.Next()
      if nxt == '/' {
        // eat every character until a newline is found
        for nxt != '\n' && nxt != '\r' {
          nxt = scanner.Next()
        }
        continue
      } else {
        // otherwise put both of the characters back and keep going
        scanner.Unread()
      }
    }

    if c == '"' {
      tkns = append(tkns, scanString(scanner))
    } else if c == '-' {
      nxt := scanner.Next()
      if numberRe.Match([]byte{nxt}) {
        scanner.Unread()
        tkns = append(tkns, scanNumber(c, scanner))
      }
    } else if isSpecial(c) {
      //tkns = append(tkns, &Token{Type: string(c)}) // special characters {}()@,;:= will be their own type
      tkns = append(tkns, token.NewToken(string(c), string(c), *scanner))
    } else if isOperation(c) {
      tkns = append(tkns, scanOperation(c, scanner))
    } else if c == '=' {
      // check the next one is '=' to see if this is special or an operation
      nxt := scanner.Next()
      if nxt == '=' {
        //tkns = append(tkns, &Token{Type: "operation", Value: "=="})
        tkns = append(tkns, token.NewToken("operation", "==", *scanner))
      } else {
        scanner.Unread()
        //tkns = append(tkns, &Token{Type: string(c)})
        tkns = append(tkns, token.NewToken(string(c), string(c), *scanner))
      }
    } else if numberRe.Match([]byte{c}) { // number literal
      tkns = append(tkns, scanNumber(c, scanner))
    } else if symbolRe.Match(([]byte{c})) { // symbol
      tkn := scan(c, scanner, symbolRe, "symbol")
      if boolRe.MatchString(tkn.Value) { // boolean
        //tkns = append(tkns, &Token{Type: "bool", Value: tkn.Value})
        tkns = append(tkns, token.NewToken("bool", tkn.Value, *scanner))
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
func scanString(scanner *token.Scanner) *token.Token {
  var val string
  c := scanner.Next()
  for c != '"' {
    val += string(c)
    c = scanner.Next()
    if c == '\\' {
      val += string(c)
      c = scanner.Next()
      val += string(c)
      c = scanner.Next()
    }
  }
  return token.NewToken("string", val, *scanner)
}

func scanOperation(c byte, scanner *token.Scanner) *token.Token {
  val := string(c)
  b := scanner.Next()

  if b == '=' {
    val += string(b)
  } else {
    scanner.Unread()
  }

  return token.NewToken("operation", val, *scanner)
}

func scanNumber(c byte, scanner *token.Scanner) *token.Token {
  var val string
  b := c
  if c == '-' {
    val += "-"
    b = scanner.Next()
  }
  for numberRe.Match([]byte{b}) {
    val += string(b)
    b = scanner.Next()
    if b == token.EOF {
      break
    }
  }
  scanner.Unread()
  return token.NewToken("number", val, *scanner)
}

func scan(c byte, scanner *token.Scanner, re *regexp.Regexp, t string) *token.Token {
  var val string
  b := c
  for re.Match([]byte{b}) {
    val += string(b)
    b = scanner.Next()
    if b == token.EOF {
      break
    }
  }
  scanner.Unread()
  return token.NewToken(t, val, *scanner)
}

