package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mcjcloud/taurine/evaluator"
	"github.com/mcjcloud/taurine/parser"
)

func main() {
  var printAst = flag.Bool("print-ast", false, "print abstract syntax tree")
  var printTokens = flag.Bool("print-tokens", false, "print the tokenized source code")
  flag.Parse()

  // check for file 
  if len(flag.Args()) < 1 {
    fmt.Println("Please provide a filename")
    os.Exit(1)
  }

  absPath, err := filepath.Abs(flag.Arg(0))
  if err != nil {
    fmt.Printf("Could not get absolute path to source file: %s\n", err.Error())
    os.Exit(1)
  }

  // create parse context
  ctx, err := parser.NewParseContext(absPath)
  if err != nil {
    fmt.Printf("Could not create parse context: %s\n", err.Error())
    os.Exit(1)
  }

  // optionally print tokens
  if *printTokens {
    ctx.CurrentIterator().PrintTokens()
    os.Exit(0)
  }

  // parse using context
  tree := parser.Parse(ctx)
  ctx.PopImportWithTree(tree)

  // TODO: check for import cycles
  if cycles := ctx.ImportGraph.FindCycles(); len(cycles) > 0 {
    fmt.Println("import cycle found.")
    for _, n := range cycles {
      fmt.Println(n)
    }
    os.Exit(1)
  }

  // print any errors during parsing
  if ctx.HasErrors() {
    ctx.PrintErrors()
    os.Exit(1)
  }

  // check for '--ast' flag
  if *printAst {
    fmt.Println(tree)
    os.Exit(0)
  }

  // evaluate 
  err = evaluator.Evaluate(tree, ctx.ImportGraph)
  if err != nil {
    fmt.Printf("eval error: %s", err)
  }
}

