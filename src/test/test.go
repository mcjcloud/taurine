package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mcjcloud/taurine/ast"
	"github.com/mcjcloud/taurine/evaluator"
	"github.com/mcjcloud/taurine/lexer"
	"github.com/mcjcloud/taurine/parser"

	"github.com/kylelemons/godebug/diff"
)

func main() {
  var path string
  if len(os.Args) < 2 {
    path = "test"
  } else {
    path = os.Args[1]
  }

  // loop through all folders in the present directory for testing
  infos, err := ioutil.ReadDir(path)
  if err != nil {
    fmt.Printf("error: %s\n", err.Error())
  }
  dirs := make([]os.FileInfo, 0)
  for _, info := range infos {
    if info.IsDir() {
      dirs = append(dirs, info)
    }
  }

  // for each directory, run the tests
  for _, dir := range dirs {
    err = testDirectory(filepath.Join(path, dir.Name()))
    if err != nil {
      fmt.Printf("error: %s\n", err.Error())
    }
    fmt.Println()
    fmt.Println()
  }
}

func testDirectory(path string) error {
  fmt.Printf("testing %s\n", path)
  // read the expected ast
  bytes, err := ioutil.ReadFile(filepath.Join(path, "ast.json"))
  if err != nil {
    return err
  }
  expectedAst := strings.TrimSpace(string(bytes))

  // read the expected program output 
  bytes, err = ioutil.ReadFile(filepath.Join(path, "output.txt"))
  if err != nil {
    return err
  }
  expectedOutput := string(bytes)

  // read the source code
  bytes, err = ioutil.ReadFile(filepath.Join(path, "src.tc"))
  if err != nil {
    return err
  }
  src := string(bytes)

  // build the AST from source
  tkns := lexer.Analyze(src)
  stmts, err := parser.Parse(tkns)
  if err != nil {
    return err
  }
  j, err := parser.JsonAst(stmts)
  if err != nil {
    return err
  }
  ast := strings.TrimSpace(j)

  // compare the AST results
  fmt.Printf("testing ast... ")
  if ast != expectedAst {
    return errors.New(fmt.Sprintf("expected ast does not match actual ast\n%s", diff.Diff(expectedAst, ast)))
  } else {
    fmt.Printf("done\n")
  }

  // evaluate test code
  fmt.Printf("testing output... ")
  if err := evaluateTestCode(path, stmts); err != nil {
    return err
  }

  // check results of execution
  outBytes, err := os.ReadFile(filepath.Join(path, "output.tmp"))
  output := string(outBytes)
  if output != expectedOutput {
    return errors.New(fmt.Sprintf("expected output does not match actual output\n%s", diff.Diff(expectedOutput, output)))
  } else {
    fmt.Printf("done")
  }

  return nil
}

func evaluateTestCode(path string, stmts *ast.BlockStatement) error {
  // set stdin to input.txt for program execution
  in, err := os.Open(filepath.Join(path, "input.txt"))
  if err != nil {
    return err
  }
  defer in.Close()
  oldStdin := os.Stdin
  defer func() { os.Stdin = oldStdin }()
  os.Stdin = in

  // set stdout to output.tmp for later comparison
  out, err := os.Create(filepath.Join(path, "output.tmp"))
  if err != nil {
    return err
  }
  defer out.Close()
  oldStdout := os.Stdout
  defer func() { os.Stdout = oldStdout }()
  os.Stdout = out

  // execute the program
  return evaluator.Evaluate(stmts)
}

