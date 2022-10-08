package llvm

import (
	"errors"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/mcjcloud/taurine/pkg/ast"
)

func (m LlvmModule) compileStatement(stmt ast.Statement, block *ir.Block) error {
	switch stmt := stmt.(type) {
	case *ast.EtchStatement:
		return m.compileEtchStatement(stmt, block)
	case *ast.ReadStatement:
		// TODO
		return errors.New("unimplemented")
	case *ast.ExpressionStatement:
		_, err := m.compileExpression(stmt.Expression, block)
		return err
	case *ast.BlockStatement:
		// TODO
		return errors.New("unimplemented")
	case *ast.IfStatement:
		// TODO
		return errors.New("unimplemented")
	case *ast.ForLoopStatement:
		// TODO
		return errors.New("unimplemented")
	case *ast.WhileLoopStatement:
		// TODO
		return errors.New("unimplemented")
	case *ast.ReturnStatement:
		// TODO
		return errors.New("unimplemented")
	default:
		return fmt.Errorf("unknown statement %s", stmt)
	}
}

func (m LlvmModule) compileEtchStatement(stmt *ast.EtchStatement, block *ir.Block) error {
	var values []value.Value
	for _, exp := range stmt.Expressions {
		e, err := m.compileExpression(exp, block)
		if err != nil {
			return err
		}
		values = append(values, e)
	}

	for _, val := range values {
		m.puts(val, block)
	}
	return nil
}
