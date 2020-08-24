package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mcjcloud/taurine/evaluator"
	"github.com/mcjcloud/taurine/lexer"
	"github.com/mcjcloud/taurine/parser"
)

// func main() {
// 	// tokens := lexer.Analyze(`first "Hello"+ 2 += myVar -=1 _another_var *= "Another\" \n string" ++ += 0.2 2005 last !=`)
// 	tokens := lexer.Analyze(`var (str) x = "hello";var (num) y = 3.14;etch 3, 2, y, x;`)
// 	// for _, tkn := range tokens {
// 	// 	fmt.Printf("%s %s\n", tkn.Type, tkn.Value)
// 	// }
// 	stmts, err := parser.Parse(tokens)
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 	} else {
// 		j, _ := json.Marshal(stmts)
// 		fmt.Printf("%s\n", string(j))
// 	}
// 	evaluator.Evaluate(stmts)
// }

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
	// j, _ := json.Marshal(stmts)
	// fmt.Printf("%s\n", string(j))
	if err != nil {
		fmt.Printf("Parsing Error: %v\n", err)
		os.Exit(1)
	}
	evaluator.Evaluate(stmts)
}
