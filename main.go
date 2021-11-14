package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mcjcloud/taurine/evaluator"
	"github.com/mcjcloud/taurine/lexer"
	"github.com/mcjcloud/taurine/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a filename")
		os.Exit(1)
	}
	bytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Could not read file %s\n", os.Args[1])
		os.Exit(1)
	}
	src := string(bytes)

	// check for '--ast' flag
	var printAst bool
	if len(os.Args) >= 3 && os.Args[2] == "--ast" {
		printAst = true
	}
	tkns := lexer.Analyze(src)
	stmts, err := parser.Parse(tkns)
	if printAst {
    j, err := parser.JsonAst(stmts)
    if err != nil {
      fmt.Printf("could not create AST JSON\n")
    }
		fmt.Printf("%s\n", j)
	}
	if err != nil {
		fmt.Printf("Parsing Error: %v\n", err)
		os.Exit(1)
	}
	err = evaluator.Evaluate(stmts)
	if err != nil {
		fmt.Printf("eval error: %s", err)
	}
}
