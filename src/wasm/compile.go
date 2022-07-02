package wasm

import (
	"errors"

	"github.com/mcjcloud/taurine/ast"
	"github.com/mcjcloud/taurine/util"
)

type Opcode byte

const (
  BLOCK Opcode = 0x02
  LOOP = 0x03

)

func Compile(tree *ast.Ast, importGraph *util.ImportGraph) error {
  // check that the ast has a blockstatement
  var block *ast.BlockStatement
  if b, ok := tree.Statement.(*ast.BlockStatement); !ok {
    return errors.New("ast must contain a block statement")
  } else {
    block = b
  }

  // compile block statements

  return nil
}
