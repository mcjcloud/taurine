package lexer

import (
  "testing"

  "github.com/mcjcloud/taurine/lexer"
  "github.com/mcjcloud/taurine/parser"
)

func TestAnalyze(t *testing.T) {
  src := "var (num) x = 0;"
  want := "{\"statements\":[{\"symbol\":\"x\",\"symbolType\":\"num\",\"value\":{\"Value\":0}}]}"
  ast, err := parser.Parse(lexer.Analyze(src))
  if err != nil {
    t.Fatalf("err: %s\n", err.Error())
  }
  if got, err := parser.JsonAst(ast); err == nil && want != got {
    t.Fatalf("wanted %s but got %s\n", want, got)
  }
}

