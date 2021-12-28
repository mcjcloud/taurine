package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "os"

  "github.com/mcjcloud/taurine/evaluator"
  "github.com/mcjcloud/taurine/lexer"
  "github.com/mcjcloud/taurine/parser"
)


func main() {
  var ast = flag.Bool("ast", false, "print abstract syntax tree")
  var printTokens = flag.Bool("print-tokens", false, "print the tokenized source code")
  flag.Parse()

  if len(flag.Args()) < 1 {
    fmt.Println("Please provide a filename")
    os.Exit(1)
  }
  bytes, err := ioutil.ReadFile(flag.Arg(0))
  if err != nil {
    fmt.Printf("Could not read file %s\n", flag.Arg(1))
    os.Exit(1)
  }
  src := string(bytes)

  tkns := lexer.Analyze(src)
  if *printTokens {
    lexer.PrintTokens(tkns)
    os.Exit(0)
  }
  it := lexer.NewTokenIterator(tkns)
  stmts := parser.Parse(it)

  // print any errors during parsing
  if len(it.EHandler.Errors) > 0 {
    fmt.Printf("found %d error(s).\n", len(it.EHandler.Errors))
    it.PrintErrors()
    os.Exit(1)
  }

  // check for '--ast' flag
  if *ast {
    j, err := parser.JsonAst(stmts)
    if err != nil {
      fmt.Printf("could not create AST JSON: %s\n", err.Error())
    }
    fmt.Printf("%s", j)
  } else {
    err = evaluator.Evaluate(stmts)
    if err != nil {
      fmt.Printf("eval error: %s", err)
    }
  }
}
