package llvm

import (
	"errors"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mcjcloud/taurine/pkg/ast"
	"github.com/mcjcloud/taurine/pkg/util"
)

type LlvmModule struct {
	*ir.Module
	globalFunctions map[string]value.Value
}

func Compile(e *ast.Ast, importGraph *util.ImportGraph) (*ir.Module, error) {
	m := LlvmModule{
		Module:          ir.NewModule(),
		globalFunctions: make(map[string]value.Value),
	}

	// setup libc bindings
	// TODO: think about replacing this with assembly to avoid C dep
	m.compileLibC()

	main := m.NewFunc("main", types.I32)
	entry := main.NewBlock("")

	// loop over block statement
	var block *ast.BlockStatement
	if b, ok := e.Statement.(*ast.BlockStatement); !ok {
		return nil, errors.New("ast must contain block statement")
	} else {
		block = b
	}

	var err error
	for _, stmt := range block.Statements {
		err = m.compileStatement(stmt, entry)
		if err != nil {
			break
		}
	}

	if err == nil {
		entry.NewRet(constant.NewInt(types.I32, 0))
	} else {
		entry.NewRet(constant.NewInt(types.I32, 1))
	}

	return m.Module, nil
}
