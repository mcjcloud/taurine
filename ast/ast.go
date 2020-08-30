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
func (b *BlockStatement) String() string {
	var str string
	for _, exp := range b.Statements {
		str += exp.String()
	}
	return str
}

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
	return fmt.Sprintf("var (%s) %s = %s", v.SymbolType, v.Symbol, v.Value)
}

// FunctionDecleration represents a function decleration
type FunctionDecleration struct {
	Symbol     string                 `json:"symbol"`
	ReturnType string                 `json:"returnType"`
	Parameters []*VariableDecleration `json:"parameters"`
	Body       Statement              `json:"body"`
}

func (f *FunctionDecleration) do() {}
func (f *FunctionDecleration) String() string {
	return fmt.Sprintf("func (%s) %s(%s) %s", f.ReturnType, f.Symbol, f.Parameters, f.Body)
}

// ReturnStatement represents a statement to return a value
type ReturnStatement struct {
	Value Expression `json:"value"`
}

func (r *ReturnStatement) do() {}
func (r *ReturnStatement) String() string {
	return fmt.Sprintf("return %s", r.Value)
}

// EtchStatement represents an etch call
type EtchStatement struct {
	Expressions []Expression `json:"expressions"`
}

func (e *EtchStatement) do() {}
func (e *EtchStatement) String() string {
	var val string
	for _, exp := range e.Expressions {
		val += exp.String()
	}
	return val
}

// ReadStatement represents a statement to read from stdin
type ReadStatement struct {
	Identifier *Identifier    `json:"expressions"`
	Prompt     *StringLiteral `json:"prompt"`
}

func (r *ReadStatement) do() {}
func (r *ReadStatement) String() string {
	return fmt.Sprintf("read %s, %s", r.Identifier, r.Prompt)
}

// IfStatement represents an if statement
type IfStatement struct {
	Condition Expression `json:"condition"`
	Statement Statement  `json:"statement"`
}

func (i *IfStatement) do() {}
func (i *IfStatement) String() string {
	return fmt.Sprintf("if %s %s", i.Condition, i.Statement)
}

// WhileLoopStatement represents a for loop
type WhileLoopStatement struct {
	Condition Expression `json:"condition"`
	Statement Statement  `json:"statement"`
}

func (w *WhileLoopStatement) do() {}
func (w *WhileLoopStatement) String() string {
	return w.Condition.String()
}

// Symbol is a type which represents the possible beginning symbols of a statement
type Symbol string

const (
	// IF represents the if keyword
	IF = "if"
	// FOR represents the for keyword
	FOR = "for"
	// WHILE represents the while keyword
	WHILE = "while"
	// VAR represents the var keyword
	VAR = "var"
	// ETCH represents the etch keyword
	ETCH = "etch"
	// READ represents the read keyword
	READ = "read"
	// RETURN represents the return keyword
	RETURN = "return"
	// NUM represents a number type
	NUM = "num"
	// STR represents a string type
	STR = "str"
	// BOOL represents a boolean type
	BOOL = "bool"
	// FUNC represents the function keyword
	FUNC = "func"
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
	return str == IF || str == FOR || str == WHILE || str == VAR || str == ETCH || str == READ || str == RETURN || str == FUNC
}

// IsDataType returns true if the symbol represents a data type
func (str Symbol) IsDataType() bool {
	return str == NUM || str == STR || str == BOOL
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

// BooleanLiteral represents a bool
type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) evaluate() {}
func (b *BooleanLiteral) String() string {
	return fmt.Sprintf("%v", b.Value)
}

// FunctionCall represents an expression which needs to call a function
type FunctionCall struct {
	Function  string       `json:"function"`
	Arguments []Expression `json:"arguments"`
}

func (f *FunctionCall) evaluate() {}
func (f *FunctionCall) String() string {
	return fmt.Sprintf("%s(%s)", f.Function, f.Arguments)
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

// AssignmentExpression represents an expression which assigns a new value to a variable
type AssignmentExpression struct {
	Identifier *Identifier `json:"identifier"`
	Value      Expression  `json:"value"`
}

func (a *AssignmentExpression) evaluate() {}
func (a *AssignmentExpression) String() string {
	return fmt.Sprintf("%s = %s", a.Identifier, a.Value)
}
