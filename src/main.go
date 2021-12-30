package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mcjcloud/taurine/evaluator"
	"github.com/mcjcloud/taurine/lexer"
	"github.com/mcjcloud/taurine/parser"
	"github.com/mcjcloud/taurine/util"
)


func main() {
  var printAst = flag.Bool("print-ast", false, "print abstract syntax tree")
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
  absPath, err := filepath.Abs(flag.Arg(0))
  it := lexer.NewTokenIterator(tkns, absPath, util.NewImportGraph())
  if err != nil {
    fmt.Printf("Could not get absolute path to source file: %s\n", err.Error())
    os.Exit(1)
  }
  tree := parser.Parse(it)

  // print any errors during parsing
  if len(it.EHandler.Errors) > 0 {
    fmt.Printf("found %d error(s).\n", len(it.EHandler.Errors))
    it.PrintErrors()
    os.Exit(1)
  }

  // check for '--ast' flag
  if *printAst {
    fmt.Println(tree)
    os.Exit(0)
  }

  // TODO: check for import cycles
  it.IGraph.Print(absPath)

  // evaluate 
  err = evaluator.Evaluate(tree, it.IGraph)
  if err != nil {
    fmt.Printf("eval error: %s", err)
  }
}

