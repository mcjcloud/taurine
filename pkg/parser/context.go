package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/mcjcloud/taurine/pkg/ast"
	"github.com/mcjcloud/taurine/pkg/lexer"
	"github.com/mcjcloud/taurine/pkg/util"
)

// ParseContext keeps track of data during parsing
type ParseContext struct {
	MainPath      string                          // the absolute file path of the main source file
	ParseStack    *util.Stack                     // the abs path of the file currently being parsed
	Iterators     map[string]*lexer.TokenIterator // the token iterators for all files
	ErrorHandlers map[string]*util.ErrorHandler   // the error handlers for all files
	ImportGraph   *util.ImportGraph               // the import gragh

	currentNode *util.ImportNode
}

func NewParseContext(absPath string) (*ParseContext, error) {
	ctx := &ParseContext{
		MainPath:      absPath,
		ParseStack:    util.NewStackWith(absPath),
		Iterators:     make(map[string]*lexer.TokenIterator),
		ErrorHandlers: make(map[string]*util.ErrorHandler),
		ImportGraph:   util.NewImportGraph(absPath),
	}

	// read source code for main file and create tokens
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("error reading referenced source: %s", err.Error())
	}
	src := string(bytes)
	tkns, err := lexer.Analyze(src)
  if err != nil {
    return nil, err
  }

	// assign token iterator and error handlers
	ctx.Iterators[absPath] = lexer.NewTokenIterator(tkns)
	ctx.ErrorHandlers[absPath] = util.NewErrorHandler()

	// setup currentNode
	ctx.currentNode = ctx.ImportGraph.Nodes[absPath]

	return ctx, nil
}

// CurrentFilePath returns the current file path
func (ctx *ParseContext) CurrentFilePath() string {
	return ctx.ParseStack.Top()
}

// CurrentFileDir returns the directory the current file is in
func (ctx *ParseContext) CurrentFileDir() string {
	return filepath.Dir(ctx.CurrentFilePath())
}

// CurrentIterator returns the iterator for the current file
func (ctx *ParseContext) CurrentIterator() *lexer.TokenIterator {
	return ctx.Iterators[ctx.ParseStack.Top()]
}

// CurrentErrorHandler returns the error handler for the current file
func (ctx *ParseContext) CurrentErrorHandler() *util.ErrorHandler {
	return ctx.ErrorHandlers[ctx.CurrentFilePath()]
}

// PushImport creates an iterator for an import in the currently iterated file
func (ctx *ParseContext) PushImport(relativePath string) error {
	// use the current path to get the absoulte path of the one being referenced
	absPath := util.ResolveImport(ctx.CurrentFileDir(), relativePath)

	// check that the import graph doesn't already contain a parsed AST for this path
	if _, ok := ctx.ImportGraph.Nodes[absPath]; ok {
		// add the connection in the import tree and indicate that it's already parsed
		ctx.ImportGraph.Add(ctx.CurrentFilePath(), absPath)
		return &util.AlreadyParsedError{
			Path: absPath,
		}
	}

	// if the path is a directory, append the directory as the file name
	absStat, err := os.Stat(absPath)
	if err != nil {
		return err
	}
	if absStat.IsDir() {
		absPath = path.Join(absPath, path.Base(absPath)) + ".tc"
	}

	// read source code for  absPath and tokenize
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("error reading referenced source: %s", err.Error())
	}
	src := string(bytes)
	tkns, err := lexer.Analyze(src)
  if err != nil {
    return fmt.Errorf("error in lexical analyzer: %s", err.Error())
  }

	// add iterator and error handler to context, push the current file
	ctx.Iterators[absPath] = lexer.NewTokenIterator(tkns)
	ctx.ErrorHandlers[absPath] = util.NewErrorHandler()

	// update import graph
	ctx.currentNode = ctx.ImportGraph.Add(ctx.ParseStack.Top(), absPath)

	// push the new path to the parsestack
	ctx.ParseStack.Push(absPath)

	return nil
}

// PopImportWithTree pops the currently parsing file, assigning the given AST to the current node
func (ctx *ParseContext) PopImportWithTree(tree *ast.Ast) {
	ctx.currentNode.SetAst(tree)
	ctx.ParseStack.Pop()
	newCurr := ctx.ParseStack.Top()
	ctx.currentNode = ctx.ImportGraph.Nodes[newCurr]
}

// HasErrors returns true if any of the error handlers are nonempty
func (ctx ParseContext) HasErrors() bool {
	for _, handler := range ctx.ErrorHandlers {
		if len(handler.Errors) > 0 {
			return true
		}
	}
	return false
}

// PrintErrors prints all errors found during parsing
func (ctx *ParseContext) PrintErrors() {
	for path, handler := range ctx.ErrorHandlers {
		if len(handler.Errors) == 0 {
			continue
		}
		it := ctx.Iterators[path]

		fmt.Printf("found %d errors in %s\n", len(handler.Errors), path)
		for _, e := range handler.Errors {
			// print error message
			fmt.Printf("%d:%d: %s\n", e.Token.Position.Row, e.Token.Position.Col, e.Message)

			// print each token in the row with the error
			row := it.GetRow(e.Token.Position.Row)
			colStart := 1
			for _, t := range row {
				// print spaces leading up to the beginning of each token
				for i := colStart; i < t.Position.Col; i += 1 {
					fmt.Printf(" ")
				}
				// update colStart and print the token
				colStart = t.Position.Col + t.Position.Length
				if t.Type == "string" {
					fmt.Printf("\"%s\"", t.Value)
				} else {
					fmt.Printf(t.Value)
				}
			}
			fmt.Println()

			// print underlines up until the error token
			for i := 0; i < e.Token.Position.Col+e.Token.Position.Length-1; i += 1 {
				fmt.Printf("~")
			}
			fmt.Println("^")
		}
	}
}
