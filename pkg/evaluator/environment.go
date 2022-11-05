package evaluator

import "github.com/mcjcloud/taurine/pkg/ast"

// Function represents a function in a particular scope
type ScopedFunction struct {
	Scope    *Scope
	Function *ast.FunctionLiteral
}

// implement Expression interface to allow scope to be stored with the function
func (s *ScopedFunction) Evaluate() {}
func (s *ScopedFunction) String() string {
	return s.Function.String()
}

// Scope represents data within a scope during execution
type Scope struct {
	Parent      *Scope                    // the parent scope
	Variables   map[string]ast.Expression // a map of variable names to values
	ReturnValue ast.Expression            // if the scope is for a function, this will hold the return value
}

// NewScope creates a new Scope
func NewScope() *Scope {
	return &Scope{
		Parent:    nil,
		Variables: map[string]ast.Expression{},
	}
}

// NewScopeWithParent creates a new scope with a parent scope
func NewScopeWithParent(par *Scope) *Scope {
	return &Scope{
		Parent:    par,
		Variables: map[string]ast.Expression{},
	}
}

// NewScopeOfObject creates a new scope with an objects properties as variables
func NewScopeOfObject(obj *ast.ObjectLiteral, par *Scope) *Scope {
	return &Scope{
		Parent:    par,
		Variables: obj.Value,
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
