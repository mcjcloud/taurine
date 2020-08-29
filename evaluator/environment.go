package evaluator

import "github.com/mcjcloud/taurine/ast"

// Scope represents data within a scope during execution
type Scope struct {
	Parent    *Scope
	Variables map[string]ast.Expression
	Functions map[string]*ast.FunctionDecleration
}

// NewScope creates a new Scope
func NewScope() *Scope {
	return &Scope{
		Parent:    nil,
		Variables: map[string]ast.Expression{},
		Functions: map[string]*ast.FunctionDecleration{},
	}
}

// NewScopeWithParent creates a new scope with a parent scope
func NewScopeWithParent(par *Scope) *Scope {
	return &Scope{
		Parent:    par,
		Variables: map[string]ast.Expression{},
	}
}

// Get returns the current value for a symbol
func (s *Scope) Get(symbol string) ast.Expression {
	if val, ok := s.Variables[symbol]; ok {
		return val
	}
	if s.Parent != nil {
		return s.Parent.Get(symbol)
	}
	return nil
}

// Set creates or updates a value in the scope
func (s *Scope) Set(symbol string, val ast.Expression) {
	if s.Variables[symbol] == nil && s.Parent != nil && s.Parent.Get(symbol) != nil {
		s.Parent.Set(symbol, val)
		return
	}
	s.Variables[symbol] = val
}

// GetFunction gets a function definition
func (s *Scope) GetFunction(symbol string) *ast.FunctionDecleration {
	if val, ok := s.Functions[symbol]; ok {
		return val
	}
	if s.Parent != nil {
		return s.Parent.GetFunction(symbol)
	}
	return nil
}

// SetFunction sets a function definition
func (s *Scope) SetFunction(symbol string, val *ast.FunctionDecleration) {
	if s.Functions[symbol] == nil && s.Parent != nil && s.Parent.GetFunction(symbol) != nil {
		s.Parent.SetFunction(symbol, val)
		return
	}
	s.Functions[symbol] = val
}
