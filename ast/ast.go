package ast

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
	Statements []*Statement `json:"statements"`
}

func (b *BlockStatement) do() {}

// Source represents a file's contents as a series of BlockStatements and Statements
type Source struct {
	BlockStatements []*BlockStatement `json:"blockStatements"`
}

func (c *Source) do() {}

// VariableDecleration represents a node that is a variable decleration
type VariableDecleration struct {
	Symbol     string      `json:"symbol"`
	SymbolType string      `json:"symbolType"`
	Value      interface{} `json:"value"`
}

func (v *VariableDecleration) do() {}
