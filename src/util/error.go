package util

import (
  "fmt"

  "github.com/mcjcloud/taurine/ast"
  "github.com/mcjcloud/taurine/token"
)

// ParseError represents an error during parsing
type ParseError struct {
  Message string
  Token *token.Token
}

// ErrorHandler keeps track of errors that occur during parsing
type ErrorHandler struct {
  Errors []ParseError
}

func NewErrorHandler() *ErrorHandler {
  return &ErrorHandler{
    Errors: make([]ParseError, 0),
  }
}

func (h *ErrorHandler) Add(tkn *token.Token, msg string) *ast.ErrorNode {
  h.Errors = append(h.Errors, ParseError{
    Message: msg,
    Token: tkn,
  })
  return &ast.ErrorNode{
    Token: tkn,
  }
}

// Specific errors that can occur
type AlreadyParsedError struct {
  Path string
}
func (a AlreadyParsedError) Error() string {
  return fmt.Sprintf("file %s has already been parsed.", a.Path)
}

