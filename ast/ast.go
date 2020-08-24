package ast

import "fmt"

// Node represents a node in the AST
type Node interface {
	String() string
}

// Statement represents a statement
type Statement interface {
	Node
	do()
}

// Expression represents an evaluatable expression
type Expression interface {
	Node
	evaluate()
}

// BlockStatement is a Statement which consists of multiple statements
type BlockStatement struct {
	Statements []Statement `json:"statements"`
}

func (b *BlockStatement) do() {}

// ExpressionStatement represents a statement which is just an expression
type ExpressionStatement struct {
	Expression Expression `json:"expression"`
}

func (e *ExpressionStatement) do() {}
func (e *ExpressionStatement) String() string {
	return e.Expression.String()
}

// VariableDecleration represents a node that is a variable decleration
type VariableDecleration struct {
	Symbol     string     `json:"symbol"`
	SymbolType string     `json:"symbolType"`
	Value      Expression `json:"value"`
}

func (v *VariableDecleration) do() {}
func (v *VariableDecleration) String() string {
	return ""
}

// EtchStatement represents an etch call
type EtchStatement struct {
	Type        string ""
	Expressions []Expression
}

func (e *EtchStatement) do() {}
func (e *EtchStatement) String() string {
	var val string
	for _, exp := range e.Expressions {
		val += exp.String()
	}
	return val
}

// ForLoopStatement represents a for loop
type ForLoopStatement struct {
	Condition Expression `json:"condition"`
	Statement Statement  `json:"statement"`
}

func (f *ForLoopStatement) do() {}
func (f *ForLoopStatement) String() string {
	return f.Condition.String()
}

// Symbol is a type which represents the possible beginning symbols of a statement
type Symbol string

const (
	// IF represents the if keyword
	IF Symbol = "if"
	// FOR represents the for keyword
	FOR = "for"
	// VAR represents the var keyword
	VAR = "var"
	// ETCH represents the etch keyword
	ETCH = "etch"
	// RETURN represents the return keyword
	RETURN = "return"
	// NUM represents a number type
	NUM = "num"
	// STR represents a string type
	STR = "str"
)

// Operator represents an operator
type Operator string

const (
	// PLUS represents +
	PLUS Operator = "+"
	// MINUS represents -
	MINUS = "-"
	// MULTIPLY represents *
	MULTIPLY = "*"
	// DIVIDE represents /
	DIVIDE = "/"
)

// IsStatementPrefix returns true if the symbol is a statement prefix
func (str Symbol) IsStatementPrefix() bool {
	return str == IF || str == FOR || str == VAR || str == ETCH || str == RETURN
}

// IsDataType returns true if the symbol represents a data type
func (str Symbol) IsDataType() bool {
	return str == NUM || str == STR
}

// NumberLiteral represents the num data type
type NumberLiteral struct {
	Value float64
}

func (n *NumberLiteral) evaluate() {}
func (n *NumberLiteral) String() string {
	return fmt.Sprintf("%f", n.Value)
}

// StringLiteral represents the str data type
type StringLiteral struct {
	Value string
}

func (s *StringLiteral) evaluate() {}
func (s *StringLiteral) String() string {
	return s.Value
}

// Identifier represents a variable or some kind of reference
type Identifier struct {
	Name string
}

func (i *Identifier) evaluate() {}
func (i *Identifier) String() string {
	return i.Name
}

// OperationExpression represents an expression consisting of an operation
type OperationExpression struct {
	Operator        Operator   `json:"operator"`
	LeftExpression  Expression `json:"leftExpression"`
	RightExpression Expression `json:"rightExpression"`
}

func (o *OperationExpression) evaluate() {}
func (o *OperationExpression) String() string {
	return fmt.Sprintf("%s %s %s", o.LeftExpression.String(), o.Operator, o.RightExpression.String())
}
