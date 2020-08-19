package main

import (
	"fmt"

	"github.com/mcjcloud/taurine/lexer"
)

func main() {
	tokens := lexer.Analyze(`first "Hello"+ 2 += myVar -=1 _another_var *= "Another\" \n string" ++ += 0.2 2005 last !=`)
	for _, tkn := range tokens {
		fmt.Printf("%s %s\n", tkn.Type, tkn.Value)
	}
}
