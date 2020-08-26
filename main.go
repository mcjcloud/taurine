package main

import (
	"encoding/json"
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
	stmts, err := parser.Parse(lexer.Analyze(src))
	j, _ := json.Marshal(stmts)
	fmt.Printf("%s\n", string(j))
	if err != nil {
		fmt.Printf("Parsing Error: %v\n", err)
		os.Exit(1)
	}
	err = evaluator.Evaluate(stmts)
	if err != nil {
		fmt.Printf("eval error: %s", err)
	}
}
