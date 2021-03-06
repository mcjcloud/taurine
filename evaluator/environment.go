package evaluator

import "github.com/mcjcloud/taurine/ast"

// Function represents a function in a particular scope
type Function struct {
	Scope       *Scope
	Decleration *ast.FunctionDecleration
}

// Scope represents data within a scope during execution
type Scope struct {
	Parent      *Scope                    // the parent scope
	Variables   map[string]ast.Expression // a map of variable names to values
	Functions   map[string]*Function      // a map of function name to declerations
	ReturnValue ast.Expression            // if the scope is for a function, this will hold the return value
}

// NewScope creates a new Scope
func NewScope() *Scope {
	return &Scope{
		Parent:    nil,
		Variables: map[string]ast.Expression{},
		Functions: map[string]*Function{},
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
func (s *Scope) GetFunction(symbol string) *Function {
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
	s.Functions[symbol] = &Function{
		Decleration: val,
		Scope:       s,
	}
}
