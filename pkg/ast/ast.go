package ast

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/mcjcloud/taurine/pkg/token"
)

// Node represents a node in the AST
type Node interface {
	fmt.Stringer
}

// Statement represents a statement
type Statement interface {
	Node
	do()
}

// Expression represents an evaluatable expression
type Expression interface {
	Node
	Evaluate()
}

// Ast represents the Abstract Syntax Tree for a file
type Ast struct {
	FilePath  string                `json:"file_path"` // the absolute path to the source code; will be used for referencing
	Exports   map[string]Expression `json:"exports"`   // used during execution to map variable names to resolved values
	Statement Statement             `json:"statement"` // the parsed AST root
	Evaluated bool                  `json:"evaluated"` // true if the AST as already been evaluated (used during execution)
}

func (a *Ast) String() string {
	j, err := json.Marshal(a.Statement)
	if err != nil {
		return fmt.Sprintf("error: %s", err.Error())
	}
	return string(j)
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
	ElseIf    Statement  `json:"else_if"` // this may just be a statement in the case of else or another IfStatement in case of else if
}

func (i *IfStatement) do() {}
func (i *IfStatement) String() string {
	return fmt.Sprintf("if %s %s else %s", i.Condition, i.Statement, i.ElseIf)
}

// ForLoopStatement represents for loop
type ForLoopStatement struct {
	Control   *Identifier `json:"control"`
	Iterator  Expression  `json:"iterator"`
	Step      int         `json:"step"`
	Statement Statement   `json:"statement"`
}

func (f *ForLoopStatement) do() {}
func (f *ForLoopStatement) String() string {
	return fmt.Sprintf("for %s in %s %s", f.Control, f.Iterator, f.Statement)
}

// WhileLoopStatement represents a while loop
type WhileLoopStatement struct {
	Condition Expression `json:"condition"`
	Statement Statement  `json:"statement"`
}

func (w *WhileLoopStatement) do() {}
func (w *WhileLoopStatement) String() string {
	return w.Condition.String()
}

// ImportStatement represents an import statement
type ImportStatement struct {
	Source  string        `json:"source"`
	Imports []*Identifier `json:"imports"`
}

func (i *ImportStatement) do() {}
func (i *ImportStatement) String() string {
	return fmt.Sprintf("import %s from %s", i.Imports, i.Source)
}

// ExportStatement represents an export statement
type ExportStatement struct {
	Identifier *Identifier `json:"identifier"`
	Value      Expression  `json:"value"`
}

func (e *ExportStatement) do() {}
func (e *ExportStatement) String() string {
	return fmt.Sprintf("export %s as %s", e.Value, e.Identifier)
}

// TODO: add for loops

// Symbol is a type which represents the possible beginning symbols of a statement
type Symbol string

const (
	// IF represents the if keyword
	IF = "if"
	// ELSE represents the else keyword
	ELSE = "else"
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
	// INT represents an integer type
	INT = "int"
	// STR represents a string type
	STR = "str"
	// BOOL represents a boolean type
	BOOL = "bool"
	// ARR represents an array type
	ARR = "arr"
	// OBJ represents the object type
	OBJ = "obj"
	// VOID represents void type
	VOID = "void"
	// FUNC represents the function keyword
	FUNC = "func"
	// IMPORT represents the import keyword
	IMPORT = "import"
	// EXPORT represents the export keyword
	EXPORT = "export"
	// AS represents as keyword
	AS = "as"
	// FROM represents from keyword
	FROM = "from"
	// IN represents in keyword
	IN = "in"
)

// Operator represents an operator
type Operator string

const (
	// PLUS represents +
	PLUS Operator = "+"
	// PLUS_EQUAL represents +=
	PLUS_EQUAL = "+="
	// MINUS represents -
	MINUS = "-"
	// MINUS_EQUAL represents -=
	MINUS_EQUAL = "-="
	// MULTIPLY represents *
	MULTIPLY = "*"
	// MULTIPLY_EQUAL represents *=
	MULTIPLY_EQUAL = "*="
	// DIVIDE represents /
	DIVIDE = "/"
	// DIVIDE_EQUAL represents /=
	DIVIDE_EQUAL = "/="
	// MODULO represents %
	MODULO = "%"
	// MODULO_EQUAL represents %=
	MODULO_EQUAL = "%="
	// EQUAL_EQUAL represents ==
	EQUAL_EQUAL = "=="
	// NOT_EQUAL represents !=
	NOT_EQUAL = "!="
	// LESS_THAN represents <
	LESS_THAN = "<"
	// LESS_EQUAL represents <=
	LESS_EQUAL = "<="
	// GREATER_THAN represents >
	GREATER_THAN = ">"
	// GREATER_EQUAL represents >=
	GREATER_EQUAL = ">="
	// AT represents @
	AT = "@"
	// DOT represents .
	DOT = "."
	// RANGE represents ..
	RANGE = ".."
)

var PRECEDENCE = map[Operator]int{
	PLUS:     1,
	MINUS:    1,
	MULTIPLY: 2,
	DIVIDE:   2,
	AT:       3,
	RANGE:    4,
	DOT:      5,
}

// IsStatementPrefix returns true if the symbol is a statement prefix
func (str Symbol) IsStatementPrefix() bool {
	return str == IF || str == FOR || str == WHILE || str == ETCH || str == READ || str == RETURN || str == IMPORT || str == EXPORT
}

// IsDataType returns true if the symbol represents a data type
func (str Symbol) IsDataType() bool {
	return str == NUM || str == INT || str == STR || str == BOOL || str == ARR || str == OBJ || str == FUNC || str == VOID
}

// ErrorNode represents an exoression that couldn't be parsed
type ErrorNode struct {
	Token *token.Token
}

func (e *ErrorNode) do()       {}
func (e *ErrorNode) Evaluate() {}
func (e *ErrorNode) String() string {
	return fmt.Sprintf("error node: %s", e.Token.Value)
}

// NumberLiteral represents the num data type
type NumberLiteral struct {
	Value float64
}

func (n *NumberLiteral) Evaluate() {}
func (n *NumberLiteral) String() string {
	return fmt.Sprintf("%f", n.Value)
}

// IntegerLiteral represents the int data type
type IntegerLiteral struct {
	Value *big.Int
}

func (i *IntegerLiteral) Evaluate() {}
func (i *IntegerLiteral) String() string {
	return i.Value.String()
}

// StringLiteral represents the str data type
type StringLiteral struct {
	Value string
}

func (s *StringLiteral) Evaluate() {}
func (s *StringLiteral) String() string {
	return s.Value
}

// BooleanLiteral represents a bool
type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) Evaluate() {}
func (b *BooleanLiteral) String() string {
	return fmt.Sprintf("%v", b.Value)
}

// ObjectLiteral represents the obj data type
type ObjectLiteral struct {
	Value map[string]Expression
}

func (o *ObjectLiteral) Evaluate() {}
func (o *ObjectLiteral) String() string {
	return fmt.Sprintf("%v", o.Value)
}

// FunctionLiteral represents a function
type FunctionLiteral struct {
	Symbol     string                 `json:"symbol"`
	ReturnType string                 `json:"returnType"`
	Parameters []*VariableDecleration `json:"parameters"`
	Body       Statement              `json:"body"`
}

func (f *FunctionLiteral) Evaluate() {}
func (f *FunctionLiteral) String() string {
	return fmt.Sprintf("func (%s) %s(%s) %s", f.ReturnType, f.Symbol, f.Parameters, f.Body)
}

// FunctionCall represents an expression which needs to call a function
// The "Expression" will be whatever in AST, but Evaluate to a evaluator.ScopedFunction during runtime
type FunctionCall struct {
	Function  Expression   `json:"function"`
	Arguments []Expression `json:"arguments"`
}

func (f *FunctionCall) Evaluate() {}
func (f *FunctionCall) String() string {
	return fmt.Sprintf("%s(%s)", f.Function, f.Arguments)
}

// VariableDecleration represents a node that is a variable decleration
type VariableDecleration struct {
	Symbol     string     `json:"symbol"`
	SymbolType string     `json:"symbolType"`
	Value      Expression `json:"value"`
}

func (v *VariableDecleration) Evaluate() {}
func (v *VariableDecleration) String() string {
	return fmt.Sprintf("var (%s) %s = %s", v.SymbolType, v.Symbol, v.Value)
}

// Identifier represents a variable or some kind of reference
type Identifier struct {
	Name string
}

func (i *Identifier) Evaluate() {}
func (i *Identifier) String() string {
	return i.Name
}

// OperationExpression represents an expression consisting of an operation
type OperationExpression struct {
	Operator        Operator   `json:"operator"`
	LeftExpression  Expression `json:"leftExpression"`
	RightExpression Expression `json:"rightExpression"`
}

func (o *OperationExpression) Evaluate() {}
func (o *OperationExpression) String() string {
	var l string
	if o.LeftExpression != nil {
		l = o.LeftExpression.String()
	}
	var r string
	if o.RightExpression != nil {
		r = o.RightExpression.String()
	}
	return fmt.Sprintf("%s(%s, %s)", o.Operator, l, r)
}

// AssignmentExpression represents an expression which assigns a new value to a variable
type AssignmentExpression struct {
	Identifier *Identifier `json:"identifier"`
	Value      Expression  `json:"value"`
}

func (a *AssignmentExpression) Evaluate() {}
func (a *AssignmentExpression) String() string {
	return fmt.Sprintf("%s = %s", a.Identifier, a.Value)
}

// GroupExpression represents an expression inside of []
type GroupExpression struct {
	Expression Expression `json:"expression"`
}

func (g *GroupExpression) Evaluate() {}
func (g *GroupExpression) String() string {
	return fmt.Sprintf("[%s]", g.Expression)
}

// ArrayExpression represents an array of expressions e.g. [exp1, exp2]
type ArrayExpression struct {
	Expressions []Expression `json:"expressions"`
}

func (a *ArrayExpression) Evaluate() {}
func (a *ArrayExpression) String() string {
	str := "["
	for i, e := range a.Expressions {
		if i == 0 {
			str += e.String()
		} else {
			str += fmt.Sprintf(", %v", e.String())
		}
	}
	str += "]"
	return str
}
