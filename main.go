package main

import (
	"encoding/json"
	"fmt"

	"github.com/mcjcloud/taurine/lexer"
	"github.com/mcjcloud/taurine/parser"
)

func main() {
	// tokens := lexer.Analyze(`first "Hello"+ 2 += myVar -=1 _another_var *= "Another\" \n string" ++ += 0.2 2005 last !=`)
	tokens := lexer.Analyze(`var (str) x = "hello";var (num) y = 3.14;etch 3, 2;`)
	// for _, tkn := range tokens {
	// 	fmt.Printf("%s %s\n", tkn.Type, tkn.Value)
	// }
	stmts, err := parser.Parse(tokens)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		j, _ := json.Marshal(stmts)
		fmt.Printf("%s\n", string(j))
	}
}
